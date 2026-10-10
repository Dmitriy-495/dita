package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dmitriy-495/dita/internal/meta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CatalogRepository реализует persistence мета-описания справочников.
type CatalogRepository struct {
	db *pgxpool.Pool
}

// NewCatalogRepository создаёт PostgreSQL repository.
func NewCatalogRepository(db *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{db: db}
}

func (r *CatalogRepository) Create(ctx context.Context, catalog meta.Catalog) (meta.Catalog, error) {
	if err := catalog.Validate(); err != nil {
		return meta.Catalog{}, err
	}

	err := r.db.QueryRow(
		ctx,
		`INSERT INTO meta_catalogs (name, prefix)
		 VALUES ($1, $2)
		 RETURNING id`,
		catalog.Name,
		catalog.Prefix,
	).Scan(&catalog.ID)
	if err != nil {
		return meta.Catalog{}, fmt.Errorf("create catalog: %w", err)
	}

	return catalog, nil
}

func (r *CatalogRepository) Get(ctx context.Context, id int64) (meta.Catalog, error) {
	var catalog meta.Catalog

	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, prefix
		 FROM meta_catalogs
		 WHERE id = $1`,
		id,
	).Scan(&catalog.ID, &catalog.Name, &catalog.Prefix)

	if errors.Is(err, pgx.ErrNoRows) {
		return meta.Catalog{}, meta.ErrCatalogNotFound
	}
	if err != nil {
		return meta.Catalog{}, fmt.Errorf("get catalog: %w", err)
	}

	return catalog, nil
}

func (r *CatalogRepository) Update(ctx context.Context, catalog meta.Catalog) error {
	if err := catalog.Validate(); err != nil {
		return err
	}

	tag, err := r.db.Exec(
		ctx,
		`UPDATE meta_catalogs
		 SET name = $1, prefix = $2
		 WHERE id = $3`,
		catalog.Name,
		catalog.Prefix,
		catalog.ID,
	)
	if err != nil {
		return fmt.Errorf("update catalog: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return meta.ErrCatalogNotFound
	}

	return nil
}
