package meta

import "context"

// CatalogService реализует операции Meta-Layer над мета-описаниями справочников.
type CatalogService struct {
	repository CatalogRepository
}

// NewCatalogService создаёт сервис Meta-Layer для справочников.
func NewCatalogService(repository CatalogRepository) *CatalogService {
	return &CatalogService{repository: repository}
}

// CreateCatalog создаёт мета-описание справочника.
func (s *CatalogService) CreateCatalog(ctx context.Context, catalog Catalog) (Catalog, error) {
	return s.repository.Create(ctx, catalog)
}

// GetCatalog получает мета-описание справочника.
func (s *CatalogService) GetCatalog(ctx context.Context, id int64) (Catalog, error) {
	return s.repository.Get(ctx, id)
}

// UpdateCatalog изменяет мета-описание справочника.
func (s *CatalogService) UpdateCatalog(ctx context.Context, catalog Catalog) error {
	return s.repository.Update(ctx, catalog)
}
