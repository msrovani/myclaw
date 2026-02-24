package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/msrovani/myclaw/internal/core"
)

// Manager controls access to physical SQLite databases per workspace.
// It caches active *DB instances and closes idle ones to limit file descriptors.
type Manager struct {
	mu     sync.RWMutex
	dbs    map[string]*DB
	cfg    Config
	ctx    context.Context
	cancel context.CancelFunc
}

// NewManager creates a database manager for multi-tenant isolation.
func NewManager(cfg Config) *Manager {
	if cfg.BaseDataDir == "" {
		cfg.BaseDataDir = "data"
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		dbs:    make(map[string]*DB),
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// getDBPath constructs the strictly isolated physical path for a workspace.
// Form: {BaseDataDir}/tenants/{uid}/workspaces/{workspace_id}/memory.db
func (m *Manager) getDBPath(tc core.TenantContext) string {
	// Security: validate path components strictly to prevent traversal
	validID := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	if !validID.MatchString(tc.UID) || !validID.MatchString(tc.WorkspaceID) {
		panic("db manager: invalid tenant context identifier format (path traversal detected)")
	}

	uid := filepath.Clean(tc.UID)
	wid := filepath.Clean(tc.WorkspaceID)

	if uid == "." || wid == "." || filepath.IsAbs(uid) || filepath.IsAbs(wid) {
		panic("db: invalid tenant context path traversal detected")
	}

	return filepath.Join(m.cfg.BaseDataDir, "tenants", uid, "workspaces", wid, "memory.db")
}

// GetDB retrieves or opens the isolated database for the given tenant context.
func (m *Manager) GetDB(ctx context.Context) (*DB, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("db manager: %w", err)
	}

	key := tc.UID + ":" + tc.WorkspaceID

	m.mu.RLock()
	db, ok := m.dbs[key]
	m.mu.RUnlock()
	if ok {
		return db, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check
	if db, ok := m.dbs[key]; ok {
		return db, nil
	}

	// Calculate isolated path and ensure directories exist
	dbPath := m.getDBPath(tc)
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, fmt.Errorf("db manager: failed to create tenant dir: %w", err)
	}

	// Open the isolated DB
	tenantCfg := m.cfg
	tenantCfg.Path = dbPath
	db, err = Open(tenantCfg)
	if err != nil {
		return nil, fmt.Errorf("db manager: open tenant db %s: %w", key, err)
	}

	// Run migrations uniquely for this physical DB
	if err := Migrate(db.Writer(), CoreMigrations()); err != nil {
		db.Close()
		return nil, fmt.Errorf("db manager: migrate tenant db %s: %w", key, err)
	}

	m.dbs[key] = db
	slog.Info("db manager: opened tenant database", "uid", tc.UID, "workspace_id", tc.WorkspaceID)
	return db, nil
}

// CloseAll shuts down all open tenant databases safely.
func (m *Manager) CloseAll() {
	m.cancel()
	m.mu.Lock()
	defer m.mu.Unlock()

	var wg sync.WaitGroup
	for key, db := range m.dbs {
		wg.Add(1)
		go func(k string, d *DB) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			d.Checkpoint(ctx)
			if err := d.Close(); err != nil {
				slog.Error("db manager: failed closing tenant db", "tenant", k, "error", err)
			}
		}(key, db)
	}
	wg.Wait()
	m.dbs = make(map[string]*DB)
}
