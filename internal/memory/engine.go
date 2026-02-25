package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/msrovani/myclaw/internal/core"
	"github.com/msrovani/myclaw/internal/db"
	"github.com/msrovani/myclaw/internal/providers"
)

// Engine is the Mem0-like persistence implementation over isolated SQLite DBs.
type Engine struct {
	dbMgr    *db.Manager
	provider providers.Provider
}

// NewEngine creates a new memory engine backed by the DB Manager and an LLM provider.
func NewEngine(mgr *db.Manager, provider providers.Provider) *Engine {
	return &Engine{
		dbMgr:    mgr,
		provider: provider,
	}
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

	// Generate embeddings automatically if provider is available
	if e.provider != nil {
		emb, err := e.provider.Embed(ctx, content)
		if err != nil {
			slog.Warn("memory engine: failed to generate embedding", "error", err, "id", mem.ID)
		} else {
			mem.Embedding = emb
		}
	}

	if mem.Metadata == nil {
		mem.Metadata = make(map[string]any)
	}
	mem.Metadata["layer"] = string(layer)

	metaBytes, err := json.Marshal(mem.Metadata)
	if err != nil {
		return mem, fmt.Errorf("marshal metadata: %w", err)
	}

	var embBytes []byte
	if len(mem.Embedding) > 0 {
		embBytes, _ = db.Float32ToBytes(mem.Embedding)
	}

	err = tenantDB.Write(ctx, func(tx *sql.Tx) error {
		q := `INSERT INTO memories (id, uid, workspace_id, content, session_id, agent_id, metadata, embedding, created_at, updated_at)
		      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(q, mem.ID, tc.UID, tc.WorkspaceID, mem.Content,
			nullableString(mem.SessionID), nullableString(mem.AgentID), string(metaBytes),
			embBytes,
			mem.CreatedAt.Format(time.RFC3339), mem.UpdatedAt.Format(time.RFC3339))
		return err
	})

	return mem, err
}

// SearchMemories performs hybrid retrieval for the tenant.
func (e *Engine) SearchMemories(ctx context.Context, query string, opts SearchOptions) ([]Memory, error) {
	_, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return nil, err
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}

	var queryEmbedding []float32
	if e.provider != nil {
		queryEmbedding, _ = e.provider.Embed(ctx, query)
	}

	// 1. Vector Search
	var vecResults []db.SearchResult
	if len(queryEmbedding) > 0 {
		vecResults, err = tenantDB.SearchVector(ctx, queryEmbedding, limit)
		if err != nil {
			slog.Error("memory engine: vector search failed", "error", err)
		}
	}

	// 2. FTS Search
	ftsResults, err := tenantDB.SearchFTS(ctx, query, limit)
	if err != nil {
		slog.Error("memory engine: fts search failed", "error", err)
	}

	// 3. Reciprocal Rank Fusion
	combinedRows := db.ReciprocalRankFusion(vecResults, ftsResults, 60)

	var results []Memory
	for _, res := range combinedRows {
		var meta map[string]any
		json.Unmarshal([]byte(res.Metadata), &meta)

		results = append(results, Memory{
			ID:       res.ID,
			Content:  res.Content,
			Metadata: meta,
			Score:    res.Score,
		})
	}

	if len(results) > limit {
		results = results[:limit]
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

	var embBytes []byte
	if e.provider != nil {
		emb, _ := e.provider.Embed(ctx, content)
		if len(emb) > 0 {
			embBytes, _ = db.Float32ToBytes(emb)
		}
	}

	metaBytes, _ := json.Marshal(metadata)
	updatedAt := time.Now().Format(time.RFC3339)

	err = tenantDB.Write(ctx, func(tx *sql.Tx) error {
		var res sql.Result
		var err error
		if len(embBytes) > 0 {
			res, err = tx.Exec("UPDATE memories SET content = ?, metadata = ?, embedding = ?, updated_at = ? WHERE id = ? AND uid = ? AND workspace_id = ?",
				content, string(metaBytes), embBytes, updatedAt, id, tc.UID, tc.WorkspaceID)
		} else {
			res, err = tx.Exec("UPDATE memories SET content = ?, metadata = ?, updated_at = ? WHERE id = ? AND uid = ? AND workspace_id = ?",
				content, string(metaBytes), updatedAt, id, tc.UID, tc.WorkspaceID)
		}
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

// CompactMemories evaluates recent short-term memories and summarizes them into a single LongTerm memory.
// This implements the "Recursive Summarization" cognitive pattern.
func (e *Engine) CompactMemories(ctx context.Context, sessionID string) (int64, error) {
	if e.provider == nil {
		return 0, fmt.Errorf("memory engine: compact requires an LLM provider")
	}

	tc, tenantDB, err := e.ensureTenant(ctx)
	if err != nil {
		return 0, err
	}

	// 1. Fetch all short-term memories for this session
	q := `SELECT content FROM memories
	      WHERE uid = ? AND workspace_id = ? AND session_id = ?
	      AND json_extract(metadata, '$.layer') = 'short_term'
	      ORDER BY created_at ASC`

	rows, err := tenantDB.ReadRows(ctx, q, tc.UID, tc.WorkspaceID, sessionID)
	if err != nil {
		return 0, fmt.Errorf("compact: fetch memories: %w", err)
	}
	defer rows.Close()

	var contents []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err == nil {
			contents = append(contents, c)
		}
	}

	if len(contents) < 2 {
		// Not enough content to summarize/compact
		return 0, nil
	}

	// 2. Build summarization prompt
	fullText := strings.Join(contents, "\n---\n")
	prompt := fmt.Sprintf("Summarize the following interaction into a single, high-density knowledge entry for long-term memory. Focus on facts, user preferences, and key events. Avoid conversational filler.\n\nCONTENT:\n%s", fullText)

	// 3. Call Provider to summarize
	resp, err := e.provider.Generate(ctx, providers.GenerateRequest{
		Messages: []providers.Message{
			{Role: providers.RoleSystem, Content: "You are a memory consolidation engine. Your task is to summarize short-term interactions into long-term knowledge."},
			{Role: providers.RoleUser, Content: prompt},
		},
		Temperature: 0.3, // Low temperature for factual summarization
	})
	if err != nil {
		return 0, fmt.Errorf("compact: llm generation failed: %w", err)
	}

	summary := strings.TrimSpace(resp.Content)
	if summary == "" {
		return 0, fmt.Errorf("compact: received empty summary from provider")
	}

	// 4. Atomic Write: Add new LongTerm memory and delete/archive short-term ones
	var compacted int64
	err = tenantDB.Write(ctx, func(tx *sql.Tx) error {
		// A. Add the summarized long-term memory
		summaryID := uuid.New().String()
		now := time.Now().Format(time.RFC3339)

		// Generate embedding for the summary too
		var embBytes []byte
		if emb, embErr := e.provider.Embed(ctx, summary); embErr == nil {
			embBytes, _ = db.Float32ToBytes(emb)
		}

		meta := map[string]any{"layer": string(LayerLongTerm), "compacted_from_session": sessionID}
		metaStr, _ := json.Marshal(meta)

		insQ := `INSERT INTO memories (id, uid, workspace_id, content, session_id, metadata, embedding, created_at, updated_at)
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := tx.Exec(insQ, summaryID, tc.UID, tc.WorkspaceID, summary, sessionID, string(metaStr), embBytes, now, now); err != nil {
			return err
		}

		// B. Delete the short-term memories that were just compacted
		delQ := `DELETE FROM memories
		         WHERE uid = ? AND workspace_id = ? AND session_id = ?
		         AND json_extract(metadata, '$.layer') = 'short_term'`
		res, err := tx.Exec(delQ, tc.UID, tc.WorkspaceID, sessionID)
		if err != nil {
			return err
		}
		compacted, _ = res.RowsAffected()
		return nil
	})

	if err == nil {
		slog.Info("memory engine: session compacted", "session_id", sessionID, "memories_reduced", compacted)
	}

	return compacted, err
}

// helper for NULLs
func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
