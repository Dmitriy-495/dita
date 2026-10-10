package http

import (
	"encoding/json"
	"net/http"

	"github.com/Dmitriy-495/dita/internal/meta"
)

// MetaHandler обслуживает HTTP API Meta-Layer.
type MetaHandler struct {
	service *meta.CatalogService
}

// NewMetaHandler создаёт HTTP handler Meta-Layer.
func NewMetaHandler(service *meta.CatalogService) *MetaHandler {
	return &MetaHandler{service: service}
}

// CreateCatalog обрабатывает создание мета-описания справочника.
func (h *MetaHandler) CreateCatalog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Name   string `json:"name"`
		Prefix string `json:"prefix"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	catalog, err := h.service.CreateCatalog(
		r.Context(),
		meta.Catalog{
			Name:   request.Name,
			Prefix: request.Prefix,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, catalogResponse{
		ID:     catalog.ID,
		Name:   catalog.Name,
		Prefix: catalog.Prefix,
	})
}
