package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/Dmitriy-495/dita/internal/meta"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCatalogRepository(t *testing.T) {
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

	repository := NewCatalogRepository(db)

	created, err := repository.Create(ctx, meta.Catalog{
		Name:   "Объекты эксплуатации",
		Prefix: "ОО",
	})
	if err != nil {
		t.Fatalf("create catalog: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(ctx, `DELETE FROM meta_catalogs WHERE id = $1`, created.ID)
		if err != nil {
			t.Errorf("cleanup catalog: %v", err)
		}
	})

	if created.ID <= 0 {
		t.Fatalf("created catalog ID = %d, want positive ID", created.ID)
	}

	if created.Name != "Объекты эксплуатации" {
		t.Fatalf("created catalog name = %q", created.Name)
	}

	if created.Prefix != "ОО" {
		t.Fatalf("created catalog prefix = %q", created.Prefix)
	}

	got, err := repository.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get catalog: %v", err)
	}

	if got != created {
		t.Fatalf("got catalog = %+v, want %+v", got, created)
	}

	updated := meta.Catalog{
		ID:     created.ID,
		Name:   "Обновлённые объекты",
		Prefix: "ОБ",
	}

	if err := repository.Update(ctx, updated); err != nil {
		t.Fatalf("update catalog: %v", err)
	}

	got, err = repository.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("get updated catalog: %v", err)
	}

	if got != updated {
		t.Fatalf("got updated catalog = %+v, want %+v", got, updated)
	}

	if err := repository.Update(ctx, meta.Catalog{
		ID:     -1,
		Name:   "Missing",
		Prefix: "MS",
	}); err != meta.ErrCatalogNotFound {
		t.Fatalf("update missing catalog error = %v, want ErrCatalogNotFound", err)
	}

	if _, err := repository.Get(ctx, -1); err != meta.ErrCatalogNotFound {
		t.Fatalf("get missing catalog error = %v, want ErrCatalogNotFound", err)
	}
}

func TestCatalogDatabaseConstraints(t *testing.T) {
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

	tests := []struct {
		name    string
		query   string
		args    []any
		wantErr bool
	}{
		{
			name:    "empty name",
			query:   `INSERT INTO meta_catalogs (name, prefix) VALUES ($1, $2)`,
			args:    []any{"", "ОО"},
			wantErr: true,
		},
		{
			name:    "empty prefix",
			query:   `INSERT INTO meta_catalogs (name, prefix) VALUES ($1, $2)`,
			args:    []any{"Объекты эксплуатации", ""},
			wantErr: true,
		},
		{
			name:    "null name",
			query:   `INSERT INTO meta_catalogs (name, prefix) VALUES ($1, $2)`,
			args:    []any{nil, "ОО"},
			wantErr: true,
		},
		{
			name:    "null prefix",
			query:   `INSERT INTO meta_catalogs (name, prefix) VALUES ($1, $2)`,
			args:    []any{"Объекты эксплуатации", nil},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := db.Exec(ctx, tt.query, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
