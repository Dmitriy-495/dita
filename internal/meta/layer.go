package meta

import "context"

// Layer — основной интерфейс Meta-Layer DITA.
// Он предоставляет платформенные операции над мета-объектами.
type Layer interface {
	CreateCatalog(ctx context.Context, catalog Catalog) (Catalog, error)
	GetCatalog(ctx context.Context, id int64) (Catalog, error)
	UpdateCatalog(ctx context.Context, catalog Catalog) error
}
