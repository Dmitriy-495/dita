package runtime

import (
	"context"
	"testing"

	"github.com/Dmitriy-495/dita/internal/meta"
)

type fakeMetaLayer struct {
	catalog meta.Catalog
	err     error
}

func (f fakeMetaLayer) CreateCatalog(context.Context, meta.Catalog) (meta.Catalog, error) {
	return meta.Catalog{}, nil
}

func (f fakeMetaLayer) GetCatalog(context.Context, int64) (meta.Catalog, error) {
	return f.catalog, f.err
}

func (f fakeMetaLayer) UpdateCatalog(context.Context, meta.Catalog) error {
	return nil
}

func TestLoadCatalog(t *testing.T) {
	layer := fakeMetaLayer{
		catalog: meta.Catalog{
			ID:     1,
			Name:   "Объекты эксплуатации",
			Prefix: "ОО",
		},
	}

	got, err := LoadCatalog(context.Background(), layer, 1)
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}

	want := Catalog{
		ID:     1,
		Name:   "Объекты эксплуатации",
		Prefix: "ОО",
	}

	if got != want {
		t.Fatalf("LoadCatalog() = %+v, want %+v", got, want)
	}
}
