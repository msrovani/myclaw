package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

// Migration represents a single database migration.
type Migration struct {
	Version     int
	Description string
	Up          string
}

// Migrate runs all pending migrations in order.
func Migrate(db *sql.DB, migrations []Migration) error {
	// Ensure migrations table exists.
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version     INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at  TEXT NOT NULL
		)
	`); err != nil {
		return fmt.Errorf("migrate: create migrations table: %w", err)
	}

	// Sort migrations by version.
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Get current version.
	var maxVersion int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&maxVersion)
	if err != nil {
		return fmt.Errorf("migrate: get current version: %w", err)
	}

	applied := 0
	for _, m := range migrations {
		if m.Version <= maxVersion {
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("migrate v%d: begin: %w", m.Version, err)
		}

		if _, err := tx.Exec(m.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate v%d (%s): %w", m.Version, m.Description, err)
		}

		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, description, applied_at) VALUES (?, ?, ?)",
			m.Version, m.Description, time.Now().UTC().Format(time.RFC3339),
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate v%d: record: %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migrate v%d: commit: %w", m.Version, err)
		}

		slog.Info("migration applied",
			"version", m.Version,
			"description", m.Description,
		)
		applied++
	}

	if applied > 0 {
		slog.Info("migrations complete", "applied", applied, "current_version", migrations[len(migrations)-1].Version)
	} else {
		slog.Info("migrations: already up to date", "version", maxVersion)
	}

	return nil
}

// CoreMigrations returns the base set of migrations for XXXCLAW.
// All tables include uid and workspace_id for defense-in-depth isolation,
// even when using DB-per-workspace physical isolation.
func CoreMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "core tables with tenant isolation",
			Up: `
				CREATE TABLE IF NOT EXISTS memories (
					id           TEXT PRIMARY KEY,
					uid          TEXT NOT NULL,
					workspace_id TEXT NOT NULL,
					content      TEXT NOT NULL,
					session_id   TEXT,
					agent_id     TEXT,
					metadata     TEXT DEFAULT '{}',
					embedding    BLOB,
					hash         TEXT,
					created_at   TEXT NOT NULL DEFAULT (datetime('now')),
					updated_at   TEXT NOT NULL DEFAULT (datetime('now'))
				);
				CREATE INDEX IF NOT EXISTS idx_memories_tenant ON memories(uid, workspace_id);
				CREATE INDEX IF NOT EXISTS idx_memories_session ON memories(session_id);
				CREATE INDEX IF NOT EXISTS idx_memories_agent ON memories(agent_id);
				CREATE INDEX IF NOT EXISTS idx_memories_created ON memories(created_at);

				CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
					content,
					metadata,
					content='memories',
					content_rowid='rowid'
				);

				CREATE TRIGGER IF NOT EXISTS memories_ai AFTER INSERT ON memories BEGIN
					INSERT INTO memories_fts(rowid, content, metadata) VALUES (new.rowid, new.content, new.metadata);
				END;
				CREATE TRIGGER IF NOT EXISTS memories_ad AFTER DELETE ON memories BEGIN
					INSERT INTO memories_fts(memories_fts, rowid, content, metadata) VALUES ('delete', old.rowid, old.content, old.metadata);
				END;
				CREATE TRIGGER IF NOT EXISTS memories_au AFTER UPDATE ON memories BEGIN
					INSERT INTO memories_fts(memories_fts, rowid, content, metadata) VALUES ('delete', old.rowid, old.content, old.metadata);
					INSERT INTO memories_fts(rowid, content, metadata) VALUES (new.rowid, new.content, new.metadata);
				END;

				CREATE TABLE IF NOT EXISTS entities (
					id           TEXT PRIMARY KEY,
					uid          TEXT NOT NULL,
					workspace_id TEXT NOT NULL,
					name         TEXT NOT NULL,
					entity_type  TEXT NOT NULL DEFAULT 'unknown',
					metadata     TEXT DEFAULT '{}',
					created_at   TEXT NOT NULL DEFAULT (datetime('now'))
				);
				CREATE INDEX IF NOT EXISTS idx_entities_tenant ON entities(uid, workspace_id);
				CREATE UNIQUE INDEX IF NOT EXISTS idx_entities_tenant_name_type ON entities(uid, workspace_id, name, entity_type);

				CREATE TABLE IF NOT EXISTS entity_relations (
					id           TEXT PRIMARY KEY,
					uid          TEXT NOT NULL,
					workspace_id TEXT NOT NULL,
					source_id    TEXT NOT NULL REFERENCES entities(id),
					target_id    TEXT NOT NULL REFERENCES entities(id),
					relation     TEXT NOT NULL,
					weight       REAL DEFAULT 1.0,
					created_at   TEXT NOT NULL DEFAULT (datetime('now')),
					UNIQUE(uid, workspace_id, source_id, target_id, relation)
				);
				CREATE INDEX IF NOT EXISTS idx_relations_tenant ON entity_relations(uid, workspace_id);
				CREATE INDEX IF NOT EXISTS idx_relations_source ON entity_relations(source_id);
				CREATE INDEX IF NOT EXISTS idx_relations_target ON entity_relations(target_id);

				CREATE TABLE IF NOT EXISTS token_usage (
					id            INTEGER PRIMARY KEY AUTOINCREMENT,
					uid           TEXT NOT NULL,
					workspace_id  TEXT NOT NULL,
					provider      TEXT NOT NULL,
					model         TEXT NOT NULL,
					input_tokens  INTEGER NOT NULL DEFAULT 0,
					output_tokens INTEGER NOT NULL DEFAULT 0,
					cost_usd      REAL DEFAULT 0,
					session_id    TEXT,
					created_at    TEXT NOT NULL DEFAULT (datetime('now'))
				);
				CREATE INDEX IF NOT EXISTS idx_tokens_tenant ON token_usage(uid, workspace_id);
				CREATE INDEX IF NOT EXISTS idx_tokens_provider ON token_usage(provider);
				CREATE INDEX IF NOT EXISTS idx_tokens_session ON token_usage(session_id);
			`,
		},
	}
}
