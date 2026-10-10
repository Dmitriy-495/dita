package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestApplyMigrations(t *testing.T) {
	dsn := os.Getenv("DITA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("DITA_TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	migrationsDir, err := filepath.Abs("../../../migrations")
	if err != nil {
		t.Fatalf("resolve migrations directory: %v", err)
	}

	if err := ApplyMigrations(ctx, db, migrationsDir); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	var count int
	if err := db.QueryRow(
		ctx,
		`SELECT count(*) FROM schema_migrations`,
	).Scan(&count); err != nil {
		t.Fatalf("count migrations: %v", err)
	}

	if count == 0 {
		t.Fatal("expected at least one applied migration")
	}

	var exists bool
	if err := db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
			  AND table_name = 'meta_catalogs'
		)`,
	).Scan(&exists); err != nil {
		t.Fatalf("check meta_catalogs: %v", err)
	}

	if !exists {
		t.Fatal("meta_catalogs table does not exist")
	}

	if err := ApplyMigrations(ctx, db, migrationsDir); err != nil {
		t.Fatalf("reapply migrations: %v", err)
	}

	var countAfter int
	if err := db.QueryRow(
		ctx,
		`SELECT count(*) FROM schema_migrations`,
	).Scan(&countAfter); err != nil {
		t.Fatalf("count migrations after reapply: %v", err)
	}

	if countAfter != count {
		t.Fatalf(
			"migration count after reapply = %d, want %d",
			countAfter,
			count,
		)
	}
}
