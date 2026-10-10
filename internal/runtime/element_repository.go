package runtime

import "context"

// ElementRepository определяет persistence-операции элементов справочника.
type ElementRepository interface {
	CreateElement(ctx context.Context, catalog Catalog, name string) (Element, error)
	GetElement(ctx context.Context, catalogID, id int64) (Element, error)
	UpdateElement(ctx context.Context, catalogID int64, element Element) error
}
