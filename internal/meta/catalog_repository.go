package meta

import "context"

// CatalogRepository определяет persistence-операции мета-описания справочника.
type CatalogRepository interface {
	Create(ctx context.Context, catalog Catalog) (Catalog, error)
	Get(ctx context.Context, id int64) (Catalog, error)
	Update(ctx context.Context, catalog Catalog) error
}
