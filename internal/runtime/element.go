package runtime

import (
	"fmt"

	"github.com/Dmitriy-495/dita/internal/meta"
)

// Element представляет элемент Runtime-справочника.
type Element struct {
	ID   int64
	Code string
	Name string
}

// NewElement создаёт элемент справочника с автоматически рассчитанным кодом.
func NewElement(catalog meta.Catalog, id int64, name string) (Element, error) {
	if id <= 0 {
		return Element{}, fmt.Errorf("element id must be positive")
	}
	if name == "" {
		return Element{}, fmt.Errorf("element name is required")
	}
	if err := catalog.Validate(); err != nil {
		return Element{}, fmt.Errorf("catalog: %w", err)
	}

	return Element{
		ID:   id,
		Code: fmt.Sprintf("%s-%06d", catalog.Prefix, id),
		Name: name,
	}, nil
}
