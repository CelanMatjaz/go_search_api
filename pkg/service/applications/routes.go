package applications

import (
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store db.ApplicationStore
}

func NewHandler(store db.ApplicationStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) AddRoutes(r chi.Router) {
	r.Route("/applications", func(r chi.Router) {
		handler := ApplicationHandler{h.store}
		r.Get("/", service.CreateGetAllHandler[types.Application](handler, sendJsonApplications))
		r.Get("/{id}", service.CreateGetSingleHandler[types.Application](handler, sendJsonApplication))
		r.Post("/", service.CreatePostHandler[types.Application](handler, sendJsonApplication))
		r.Put("/", service.CreatePutHandler[types.Application](handler, sendJsonApplication))
		r.Delete("/{id}", service.CreateDeleteHandler[types.Application](handler, sendJsonApplication))
	})

	r.Route("/application-sections", func(r chi.Router) {
		handler := SectionHandler{h.store}
		r.Get("/", service.CreateGetAllHandler[types.ApplicationSection](handler, sendJsonSections))
		r.Post("/", service.CreatePostHandler[types.ApplicationSection](handler, sendJsonSection))
		r.Put("/", service.CreatePutHandler[types.ApplicationSection](handler, sendJsonSection))
		r.Delete("/{id}", service.CreateDeleteHandler[types.ApplicationSection](handler, sendJsonSection))
	})
}
