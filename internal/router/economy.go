package router

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/msrovani/myclaw/internal/core"
	"github.com/msrovani/myclaw/internal/db"
)

// Economy tracks and records costs and tokens consumed by tenants.
type Economy struct {
	dbMgr *db.Manager
}

// NewEconomy creates a Token Economy tracker.
func NewEconomy(mgr *db.Manager) *Economy {
	return &Economy{dbMgr: mgr}
}

// RecordUsage writes token consumption to the isolated tenant DB.
func (e *Economy) RecordUsage(ctx context.Context, provider string, model string, inTokens, outTokens int, cost float64) error {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return fmt.Errorf("token economy: %w", err)
	}

	tenantDB, err := e.dbMgr.GetDB(ctx)
	if err != nil {
		return fmt.Errorf("token economy: open db: %w", err)
	}

	return tenantDB.Write(ctx, func(tx *sql.Tx) error {
		q := `INSERT INTO token_usage (uid, workspace_id, provider, model, input_tokens, output_tokens, cost_usd, session_id) 
		      VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := tx.Exec(q, tc.UID, tc.WorkspaceID, provider, model, inTokens, outTokens, cost, nullableString(tc.SessionID))
		return err
	})
}

// GetTotalCost retrieves the total computed cost for the workspace.
func (e *Economy) GetTotalCost(ctx context.Context) (float64, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return 0, err
	}

	tenantDB, err := e.dbMgr.GetDB(ctx)
	if err != nil {
		return 0, err
	}

	var total float64
	q := "SELECT COALESCE(SUM(cost_usd), 0) FROM token_usage WHERE uid = ? AND workspace_id = ?"
	err = tenantDB.ReadRow(ctx, q, tc.UID, tc.WorkspaceID).Scan(&total)
	return total, err
}

func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
