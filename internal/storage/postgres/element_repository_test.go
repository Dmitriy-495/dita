package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Dmitriy-495/dita/internal/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestElementRepository(t *testing.T) {
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

	repository := NewElementRepository(db)

	var catalogID int64
	err = db.QueryRow(
		ctx,
		`INSERT INTO meta_catalogs (name, prefix)
		 VALUES ('Тестовый справочник', 'ТС')
		 RETURNING id`,
	).Scan(&catalogID)
	if err != nil {
		t.Fatalf("create test catalog: %v", err)
	}
	defer db.Exec(ctx, `DELETE FROM meta_catalogs WHERE id = $1`, catalogID)

	element, err := repository.CreateElement(
		ctx,
		runtime.Catalog{
			ID:     catalogID,
			Name:   "Тестовый справочник",
			Prefix: "ТС",
		},
		"Тестовый элемент",
	)
	if err != nil {
		t.Fatalf("CreateElement() error = %v", err)
	}

	if element.ID <= 0 {
		t.Fatalf("element ID = %d, want positive", element.ID)
	}

	if element.Code != "ТС-000001" {
		t.Fatalf("element Code = %q, want %q", element.Code, "ТС-000001")
	}

	got, err := repository.GetElement(ctx, catalogID, element.ID)
	if err != nil {
		t.Fatalf("GetElement() error = %v", err)
	}

	if got != element {
		t.Fatalf("GetElement() = %+v, want %+v", got, element)
	}

	element.Name = "Изменённый элемент"

	if err := repository.UpdateElement(ctx, catalogID, element); err != nil {
		t.Fatalf("UpdateElement() error = %v", err)
	}

	got, err = repository.GetElement(ctx, catalogID, element.ID)
	if err != nil {
		t.Fatalf("GetElement() after update error = %v", err)
	}

	if got.Name != "Изменённый элемент" {
		t.Fatalf("updated name = %q, want %q", got.Name, "Изменённый элемент")
	}

	if got.Code != "ТС-000001" {
		t.Fatalf("updated code = %q, want %q", got.Code, "ТС-000001")
	}

	_, err = repository.GetElement(ctx, catalogID, 999999999)
	if !errors.Is(err, runtime.ErrElementNotFound) {
		t.Fatalf("missing element error = %v, want ErrElementNotFound", err)
	}
}
