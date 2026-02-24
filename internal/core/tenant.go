package core

import (
	"context"
	"errors"
)

// TenantContext defines the mandatory context for all operations in XXXCLAW.
// The system operates on a "deny by default" access model and utilizes
// physical DB-per-workspace isolation.
type TenantContext struct {
	UID         string
	WorkspaceID string
	AgentID     string   // Optional, depending on the operational scope
	SessionID   string   // Optional
	AuthClaims  []string // Roles or permissions
	RequestID   string
}

// Ensure TenantContext implements basic validation
func (tc TenantContext) Validate() error {
	if tc.UID == "" {
		return errors.New("TenantContext: UID is required")
	}
	if tc.WorkspaceID == "" {
		return errors.New("TenantContext: WorkspaceID is required")
	}
	return nil
}

type tenantCtxKey struct{}

// WithTenant injects the TenantContext into a standard context.Context.
func WithTenant(ctx context.Context, tc TenantContext) context.Context {
	return context.WithValue(ctx, tenantCtxKey{}, tc)
}

// TenantFromContext extracts the TenantContext from a context.Context.
// Returns an error if the TenantContext is missing or invalid.
func TenantFromContext(ctx context.Context) (TenantContext, error) {
	tc, ok := ctx.Value(tenantCtxKey{}).(TenantContext)
	if !ok {
		return TenantContext{}, errors.New("missing TenantContext in context")
	}
	if err := tc.Validate(); err != nil {
		return TenantContext{}, err
	}
	return tc, nil
}
