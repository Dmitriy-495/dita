package postgres

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ApplyMigrations применяет SQL-миграции из указанного каталога по порядку.
func ApplyMigrations(ctx context.Context, db *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		files = append(files, entry.Name())
	}

	sort.Strings(files)

	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	for _, name := range files {
		if err := applyMigration(ctx, db, filepath.Join(dir, name), name); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	return nil
}

func applyMigration(ctx context.Context, db *pgxpool.Pool, path, version string) error {
	var applied bool

	err := db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM schema_migrations WHERE version = $1
		)`,
		version,
	).Scan(&applied)
	if err != nil {
		return fmt.Errorf("check migration %q: %w", version, err)
	}

	if applied {
		return nil
	}

	sql, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration %q: %w", version, err)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration %q: %w", version, err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		return fmt.Errorf("apply migration %q: %w", version, err)
	}

	if _, err := tx.Exec(
		ctx,
		`INSERT INTO schema_migrations (version) VALUES ($1)`,
		version,
	); err != nil {
		return fmt.Errorf("record migration %q: %w", version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %q: %w", version, err)
	}

	return nil
}
