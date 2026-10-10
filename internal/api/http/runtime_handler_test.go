package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dmitriy-495/dita/internal/meta"
	"github.com/Dmitriy-495/dita/internal/runtime"
)

type handlerLayer struct {
	catalog meta.Catalog
}

func (f *handlerLayer) CreateCatalog(context.Context, meta.Catalog) (meta.Catalog, error) {
	panic("unexpected call")
}

func (f *handlerLayer) GetCatalog(context.Context, int64) (meta.Catalog, error) {
	return f.catalog, nil
}

func (f *handlerLayer) UpdateCatalog(context.Context, meta.Catalog) error {
	panic("unexpected call")
}

type handlerElements struct {
	element runtime.Element
}

func (f *handlerElements) CreateElement(
	context.Context,
	runtime.Catalog,
	string,
) (runtime.Element, error) {
	return f.element, nil
}

func (f *handlerElements) GetElement(
	context.Context,
	int64,
	int64,
) (runtime.Element, error) {
	return f.element, nil
}

func (f *handlerElements) UpdateElement(
	context.Context,
	int64,
	runtime.Element,
) error {
	return nil
}

func TestRuntimeHandlerCreateElement(t *testing.T) {
	service := runtime.NewCatalogService(
		&handlerLayer{
			catalog: meta.Catalog{
				ID:     10,
				Name:   "Объекты эксплуатации",
				Prefix: "ОО",
			},
		},
		&handlerElements{
			element: runtime.Element{
				ID:   1,
				Code: "ОО-000001",
				Name: "Котельная №1",
			},
		},
	)

	handler := NewRuntimeHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/catalogs/10/elements",
		strings.NewReader(`{"name":"Котельная №1"}`),
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	if !strings.Contains(rec.Body.String(), `"code":"ОО-000001"`) {
		t.Fatalf("response = %s, want element code", rec.Body.String())
	}
}

func TestRuntimeHandlerGetElement(t *testing.T) {
	service := runtime.NewCatalogService(
		&handlerLayer{},
		&handlerElements{
			element: runtime.Element{
				ID:   1,
				Code: "ОО-000001",
				Name: "Котельная №1",
			},
		},
	)

	handler := NewRuntimeHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/catalogs/10/elements/1",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), `"name":"Котельная №1"`) {
		t.Fatalf("response = %s, want element name", rec.Body.String())
	}
}

func TestRuntimeHandlerInvalidCatalogID(t *testing.T) {
	handler := NewRuntimeHandler(nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/catalogs/not-an-id",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
