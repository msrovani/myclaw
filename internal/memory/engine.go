package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/msrovani/myclaw/internal/core"
	"github.com/msrovani/myclaw/internal/db"
)

// Engine is the Mem0-like persistence implementation over isolated SQLite DBs.
type Engine struct {
	dbMgr *db.Manager
}

// NewEngine creates a new memory engine backed by the DB Manager.
func NewEngine(mgr *db.Manager) *Engine {
	return &Engine{dbMgr: mgr}
}

// ensureTenant forces extraction of the context, returning the physical *DB for that workspace.
func (e *Engine) ensureTenant(ctx context.Context) (core.TenantContext, *db.DB, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return tc, nil, fmt.Errorf("memory engine: %w", err)
	}
	tenantDB, err := e.dbMgr.GetDB(ctx)
	if err != nil {
		return tc, nil, fmt.Errorf("memory engine: failed to get tenant db: %w", err)
	}
	return tc, tenantDB, nil
}

// AddMemory stores new knowledge into the tenant's isolated storage.
func (e *Engine) AddMemory(ctx context.Context, content string, opts AddOptions) (Memory, error) {
	tc, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return Memory{}, err
	}

	layer := opts.Layer
	if layer == "" {
		layer = LayerLongTerm
	}

	mem := Memory{
		ID:        uuid.New().String(),
		Content:   content,
		Layer:     layer,
		SessionID: opts.SessionID,
		AgentID:   opts.AgentID,
		Metadata:  opts.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if mem.Metadata == nil {
		mem.Metadata = make(map[string]any)
	}
	mem.Metadata["layer"] = string(layer)

	metaBytes, err := json.Marshal(mem.Metadata)
	if err != nil {
		return mem, fmt.Errorf("marshal metadata: %w", err)
	}

	err = tenantDB.Write(ctx, func(tx *sql.Tx) error {
		q := `INSERT INTO memories (id, uid, workspace_id, content, session_id, agent_id, metadata, created_at, updated_at) 
		      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(q, mem.ID, tc.UID, tc.WorkspaceID, mem.Content,
			nullableString(mem.SessionID), nullableString(mem.AgentID), string(metaBytes),
			mem.CreatedAt.Format(time.RFC3339), mem.UpdatedAt.Format(time.RFC3339))
		return err
	})

	return mem, err
}

// SearchMemories performs hybrid retrieval for the tenant.
func (e *Engine) SearchMemories(ctx context.Context, query string, queryEmbedding []float32, opts SearchOptions) ([]Memory, error) {
	// Full hybrid logic and rank fusion goes here in a future step.
	// For now we implement basic fallback vector retrieval to satisfy the interface.
	_, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return nil, err
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}

	vecResults, err := tenantDB.SearchVectorFallback(ctx, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	var results []Memory
	for _, res := range vecResults {
		var meta map[string]any
		json.Unmarshal([]byte(res.Metadata), &meta)

		results = append(results, Memory{
			ID:       res.ID,
			Content:  res.Content,
			Metadata: meta,
			Score:    1.0 - res.Distance, // Convert distance to similarity score
		})
	}

	return results, nil
}

// GetMemory retrieves a specific memory by ID.
func (e *Engine) GetMemory(ctx context.Context, id string) (Memory, error) {
	tc, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return Memory{}, err
	}

	q := "SELECT content, session_id, agent_id, metadata, created_at, updated_at FROM memories WHERE id = ? AND uid = ? AND workspace_id = ?"
	row := tenantDB.ReadRow(ctx, q, id, tc.UID, tc.WorkspaceID)

	var mem Memory
	mem.ID = id
	var sessionID, agentID sql.NullString
	var metaStr, createdStr, updatedStr string

	if err := row.Scan(&mem.Content, &sessionID, &agentID, &metaStr, &createdStr, &updatedStr); err != nil {
		return Memory{}, fmt.Errorf("get memory %s: %w", id, err)
	}

	mem.SessionID = sessionID.String
	mem.AgentID = agentID.String
	json.Unmarshal([]byte(metaStr), &mem.Metadata)

	if lay, ok := mem.Metadata["layer"].(string); ok {
		mem.Layer = Layer(lay)
	}

	mem.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	mem.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)

	return mem, nil
}

// UpdateMemory modifies an existing memory.
func (e *Engine) UpdateMemory(ctx context.Context, id string, content string, metadata map[string]any) (Memory, error) {
	tc, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return Memory{}, err
	}

	metaBytes, _ := json.Marshal(metadata)
	updatedAt := time.Now().Format(time.RFC3339)

	err = tenantDB.Write(ctx, func(tx *sql.Tx) error {
		res, err := tx.Exec("UPDATE memories SET content = ?, metadata = ?, updated_at = ? WHERE id = ? AND uid = ? AND workspace_id = ?",
			content, string(metaBytes), updatedAt, id, tc.UID, tc.WorkspaceID)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err != nil {
		return Memory{}, fmt.Errorf("update memory: %w", err)
	}

	return e.GetMemory(ctx, id)
}

// DeleteMemory removes a memory.
func (e *Engine) DeleteMemory(ctx context.Context, id string) error {
	tc, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return err
	}

	return tenantDB.Write(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM memories WHERE id = ? AND uid = ? AND workspace_id = ?", id, tc.UID, tc.WorkspaceID)
		return err
	})
}

// helper for NULLs
func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
