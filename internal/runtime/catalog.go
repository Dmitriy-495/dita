package runtime

import (
	"context"
	"fmt"

	"github.com/Dmitriy-495/dita/internal/meta"
)

// Catalog представляет загруженное Runtime-описание справочника.
type Catalog struct {
	ID     int64
	Name   string
	Prefix string
}

// LoadCatalog загружает Runtime-описание справочника через Meta-Layer.
func LoadCatalog(ctx context.Context, layer meta.Layer, id int64) (Catalog, error) {
	definition, err := layer.GetCatalog(ctx, id)
	if err != nil {
		return Catalog{}, fmt.Errorf("load catalog %d: %w", id, err)
	}

	return Catalog{
		ID:     definition.ID,
		Name:   definition.Name,
		Prefix: definition.Prefix,
	}, nil
}
