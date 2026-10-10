package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dmitriy-495/dita/internal/meta"
	"github.com/Dmitriy-495/dita/internal/runtime"
)

// RuntimeHandler обслуживает HTTP API Runtime.
type RuntimeHandler struct {
	service *runtime.CatalogService
}

// NewRuntimeHandler создаёт HTTP handler Runtime.
func NewRuntimeHandler(service *runtime.CatalogService) *RuntimeHandler {
	return &RuntimeHandler{service: service}
}

// ServeHTTP обрабатывает запросы Runtime API.
func (h *RuntimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 4 || parts[0] != "api" ||
		parts[1] != "v1" || parts[2] != "catalogs" {
		http.NotFound(w, r)
		return
	}

	catalogID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil || catalogID <= 0 {
		http.Error(w, "invalid catalog id", http.StatusBadRequest)
		return
	}

	if len(parts) == 4 {
		h.handleCatalog(w, r, catalogID)
		return
	}

	if parts[4] != "elements" {
		http.NotFound(w, r)
		return
	}

	switch len(parts) {
	case 5:
		h.handleElements(w, r, catalogID)
	case 6:
		elementID, err := strconv.ParseInt(parts[5], 10, 64)
		if err != nil || elementID <= 0 {
			http.Error(w, "invalid element id", http.StatusBadRequest)
			return
		}
		h.handleElement(w, r, catalogID, elementID)
	default:
		http.NotFound(w, r)
	}
}

func (h *RuntimeHandler) handleCatalog(w http.ResponseWriter, r *http.Request, catalogID int64) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	catalog, err := h.service.LoadCatalog(r.Context(), catalogID)
	if err != nil {
		if errors.Is(err, meta.ErrCatalogNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, newCatalogResponse(catalog))
}

func (h *RuntimeHandler) handleElements(w http.ResponseWriter, r *http.Request, catalogID int64) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	element, err := h.service.CreateElement(r.Context(), catalogID, request.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, newElementResponse(element))
}

func (h *RuntimeHandler) handleElement(w http.ResponseWriter, r *http.Request, catalogID, elementID int64) {
	switch r.Method {
	case http.MethodGet:
		element, err := h.service.GetElement(r.Context(), catalogID, elementID)
		if err != nil {
			if errors.Is(err, runtime.ErrElementNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, newElementResponse(element))

	case http.MethodPatch:
		var request struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		err := h.service.UpdateElement(r.Context(), catalogID, runtime.Element{
			ID: elementID, Name: request.Name,
		})
		if err != nil {
			if errors.Is(err, runtime.ErrElementNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		element, err := h.service.GetElement(r.Context(), catalogID, elementID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, newElementResponse(element))

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type catalogResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

type elementResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func newCatalogResponse(catalog runtime.Catalog) catalogResponse {
	return catalogResponse{ID: catalog.ID, Name: catalog.Name, Prefix: catalog.Prefix}
}

func newElementResponse(element runtime.Element) elementResponse {
	return elementResponse{ID: element.ID, Code: element.Code, Name: element.Name}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
