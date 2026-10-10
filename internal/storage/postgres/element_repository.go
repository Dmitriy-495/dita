package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dmitriy-495/dita/internal/meta"
	"github.com/Dmitriy-495/dita/internal/runtime"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ElementRepository реализует persistence элементов справочника.
type ElementRepository struct {
	db *pgxpool.Pool
}

// NewElementRepository создаёт PostgreSQL repository.
func NewElementRepository(db *pgxpool.Pool) *ElementRepository {
	return &ElementRepository{db: db}
}

// CreateElement создаёт элемент и формирует его код после получения id.
func (r *ElementRepository) CreateElement(
	ctx context.Context,
	catalog runtime.Catalog,
	name string,
) (runtime.Element, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return runtime.Element{}, fmt.Errorf("begin create element: %w", err)
	}
	defer tx.Rollback(ctx)

	var id int64

	err = tx.QueryRow(
		ctx,
		`INSERT INTO meta_catalog_elements (catalog_id, code, name)
		 VALUES ($1, '', $2)
		 RETURNING id`,
		catalog.ID,
		name,
	).Scan(&id)
	if err != nil {
		return runtime.Element{}, fmt.Errorf("create element: %w", err)
	}

	element, err := runtime.NewElement(
		meta.Catalog{
			ID:     catalog.ID,
			Name:   catalog.Name,
			Prefix: catalog.Prefix,
		},
		id,
		name,
	)
	if err != nil {
		return runtime.Element{}, err
	}

	if _, err := tx.Exec(
		ctx,
		`UPDATE meta_catalog_elements
		 SET code = $1
		 WHERE id = $2`,
		element.Code,
		element.ID,
	); err != nil {
		return runtime.Element{}, fmt.Errorf("set element code: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return runtime.Element{}, fmt.Errorf("commit create element: %w", err)
	}

	return element, nil
}

// GetElement получает активный элемент справочника.
func (r *ElementRepository) GetElement(
	ctx context.Context,
	catalogID, id int64,
) (runtime.Element, error) {
	var element runtime.Element

	err := r.db.QueryRow(
		ctx,
		`SELECT id, code, name
		 FROM meta_catalog_elements
		 WHERE catalog_id = $1
		   AND id = $2`,
		catalogID,
		id,
	).Scan(&element.ID, &element.Code, &element.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return runtime.Element{}, runtime.ErrElementNotFound
	}
	if err != nil {
		return runtime.Element{}, fmt.Errorf("get element: %w", err)
	}

	return element, nil
}

// UpdateElement изменяет имя элемента. Код остаётся неизменным.
func (r *ElementRepository) UpdateElement(
	ctx context.Context,
	catalogID int64,
	element runtime.Element,
) error {
	tag, err := r.db.Exec(
		ctx,
		`UPDATE meta_catalog_elements
		 SET name = $1
		 WHERE catalog_id = $2
		   AND id = $3`,
		element.Name,
		catalogID,
		element.ID,
	)
	if err != nil {
		return fmt.Errorf("update element: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return runtime.ErrElementNotFound
	}

	return nil
}
