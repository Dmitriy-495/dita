package http

import "net/http"

// Handler объединяет HTTP API Meta-Layer и Runtime.
type Handler struct {
	meta    *MetaHandler
	runtime *RuntimeHandler
}

// NewHandler создаёт общий HTTP handler DITA.
func NewHandler(
	metaHandler *MetaHandler,
	runtimeHandler *RuntimeHandler,
) *Handler {
	return &Handler{
		meta:    metaHandler,
		runtime: runtimeHandler,
	}
}

// ServeHTTP маршрутизирует запросы DITA API.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/catalogs" {
		h.meta.CreateCatalog(w, r)
		return
	}

	h.runtime.ServeHTTP(w, r)
}
