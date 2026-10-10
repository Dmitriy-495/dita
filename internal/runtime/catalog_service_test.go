package runtime

import (
	"context"
	"errors"
	"testing"

	"github.com/Dmitriy-495/dita/internal/meta"
)

type catalogServiceLayer struct {
	catalog meta.Catalog
	err     error
}

func (f *catalogServiceLayer) CreateCatalog(context.Context, meta.Catalog) (meta.Catalog, error) {
	panic("unexpected call")
}

func (f *catalogServiceLayer) GetCatalog(context.Context, int64) (meta.Catalog, error) {
	return f.catalog, f.err
}

func (f *catalogServiceLayer) UpdateCatalog(context.Context, meta.Catalog) error {
	panic("unexpected call")
}

type catalogServiceElements struct {
	catalog Catalog
	name    string
	element Element
}

func (f *catalogServiceElements) CreateElement(
	_ context.Context,
	catalog Catalog,
	name string,
) (Element, error) {
	f.catalog = catalog
	f.name = name

	return f.element, nil
}

func (f *catalogServiceElements) GetElement(context.Context, int64, int64) (Element, error) {
	panic("unexpected call")
}

func (f *catalogServiceElements) UpdateElement(context.Context, int64, Element) error {
	panic("unexpected call")
}

func TestCatalogServiceCreateElement(t *testing.T) {
	elements := &catalogServiceElements{
		element: Element{
			ID:   1,
			Code: "ОО-000001",
			Name: "Котельная №1",
		},
	}

	service := NewCatalogService(
		&catalogServiceLayer{
			catalog: meta.Catalog{
				ID:     10,
				Name:   "Объекты эксплуатации",
				Prefix: "ОО",
			},
		},
		elements,
	)

	got, err := service.CreateElement(
		context.Background(),
		10,
		"Котельная №1",
	)
	if err != nil {
		t.Fatalf("CreateElement() error = %v", err)
	}

	if got.Code != "ОО-000001" {
		t.Fatalf("CreateElement() code = %q, want %q", got.Code, "ОО-000001")
	}

	if elements.catalog.Name != "Объекты эксплуатации" {
		t.Fatalf(
			"repository received catalog name = %q, want %q",
			elements.catalog.Name,
			"Объекты эксплуатации",
		)
	}

	if elements.catalog.Prefix != "ОО" {
		t.Fatalf(
			"repository received catalog prefix = %q, want %q",
			elements.catalog.Prefix,
			"ОО",
		)
	}

	if elements.name != "Котельная №1" {
		t.Fatalf(
			"repository received element name = %q, want %q",
			elements.name,
			"Котельная №1",
		)
	}
}

func TestCatalogServiceCreateElementPropagatesCatalogError(t *testing.T) {
	wantErr := errors.New("catalog unavailable")

	service := NewCatalogService(
		&catalogServiceLayer{err: wantErr},
		&catalogServiceElements{},
	)

	_, err := service.CreateElement(
		context.Background(),
		10,
		"Котельная №1",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("CreateElement() error = %v, want %v", err, wantErr)
	}
}
