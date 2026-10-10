package runtime

import (
	"context"

	"github.com/Dmitriy-495/dita/internal/meta"
)

// CatalogService связывает Runtime-справочник с Meta-Layer и хранилищем элементов.
type CatalogService struct {
	layer    meta.Layer
	elements ElementRepository
}

// NewCatalogService создаёт Runtime-сервис справочников.
func NewCatalogService(layer meta.Layer, elements ElementRepository) *CatalogService {
	return &CatalogService{
		layer:    layer,
		elements: elements,
	}
}

// LoadCatalog загружает Runtime-описание справочника.
func (s *CatalogService) LoadCatalog(ctx context.Context, id int64) (Catalog, error) {
	return LoadCatalog(ctx, s.layer, id)
}

// CreateElement создаёт элемент Runtime-справочника.
func (s *CatalogService) CreateElement(
	ctx context.Context,
	catalogID int64,
	name string,
) (Element, error) {
	catalog, err := s.LoadCatalog(ctx, catalogID)
	if err != nil {
		return Element{}, err
	}

	return s.elements.CreateElement(ctx, catalog, name)
}

// GetElement получает элемент Runtime-справочника.
func (s *CatalogService) GetElement(
	ctx context.Context,
	catalogID, id int64,
) (Element, error) {
	return s.elements.GetElement(ctx, catalogID, id)
}

// UpdateElement изменяет элемент Runtime-справочника.
func (s *CatalogService) UpdateElement(
	ctx context.Context,
	catalogID int64,
	element Element,
) error {
	return s.elements.UpdateElement(ctx, catalogID, element)
}
