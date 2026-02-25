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
	Optional    bool // If true, failure to apply won't stop the system (e.g., missing extensions)
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
			if m.Optional {
				slog.Warn("optional migration failed (likely missing extension)", "version", m.Version, "error", err)
				continue
			}
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

	return nil
}

// CoreMigrations returns the base set of migrations for XXXCLAW.
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
			`,
		},
		{
			Version:     2,
			Description: "sqlite-vec integration",
			Optional:    true, // Permite rodar sem a extensão C sqlite-vec
			Up: `
				CREATE VIRTUAL TABLE memories_vec0 USING vec0(
					embedding float[384]
				);

				CREATE TRIGGER memories_vec_ai AFTER INSERT ON memories WHEN new.embedding IS NOT NULL BEGIN
					INSERT INTO memories_vec0(rowid, embedding) VALUES (new.rowid, new.embedding);
				END;
			`,
		},
	}
}
