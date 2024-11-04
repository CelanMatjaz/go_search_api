package handlers

import (
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type ApplicationHandler struct {
	store db.ApplicationStore
}

func NewApplicationHandler(store db.ApplicationStore) *ApplicationHandler {
	return &ApplicationHandler{store: store}
}

func (h *ApplicationHandler) AddRoutes(r chi.Router) {
	r.Route("/applications", func(r chi.Router) {
		r.Route("/presets", func(r chi.Router) {
			r.Get("/", CreateHandler(CreateGenericGetManyWithPaginationHandler(h.store.GetApplicationPresets, sendApplicationPresets)))
			r.Get("/{id}", CreateHandler(CreateGenericGetSingleHandler(h.store.GetApplicationPreset, sendApplicationPreset)))
			r.Post("/", CreateHandler(CreateGenericPostHandler(h.store.CreateApplicationPreset, sendApplicationPreset)))
			r.Put("/{id}", CreateHandler(CreateGenericPutHandler(h.store.UpdateApplicationPreset, sendApplicationPreset)))
			r.Delete("/{id}", CreateHandler(CreateGenericDeleteHandler(h.store.DeleteApplicationPreset)))
		})

		r.Route("/sections", func(r chi.Router) {
			r.Get("/", CreateHandler(CreateGenericGetManyWithPaginationHandler(h.store.GetApplicationSections, sendApplicationSections)))
			r.Get("/{id}", CreateHandler(CreateGenericGetSingleHandler(h.store.GetApplicationSection, sendApplicationSection)))
			r.Post("/", CreateHandler(CreateGenericPostHandler(h.store.CreateApplicationSection, sendApplicationSection)))
			r.Put("/{id}", CreateHandler(CreateGenericPutHandler(h.store.UpdateApplicationSection, sendApplicationSection)))
			r.Delete("/{id}", CreateHandler(CreateGenericDeleteHandler(h.store.DeleteApplicationSection)))
		})
	})
}

func sendApplicationPresets(w http.ResponseWriter, data []AppPreWithTags) error {
	return utils.SendJson(w, struct {
		Presets []AppPreWithTags `json:"applicationPresets"`
	}{Presets: data}, http.StatusOK)
}

func sendApplicationPreset(w http.ResponseWriter, data types.ApplicationPreset) error {
	return utils.SendJson(w, struct {
		Preset types.ApplicationPreset `json:"applicationPreset"`
	}{Preset: data}, http.StatusOK)
}

func sendApplicationSections(w http.ResponseWriter, data []AppSecWithTags) error {
	return utils.SendJson(w, struct {
		Sections []AppSecWithTags `json:"applicationSections"`
	}{Sections: data}, http.StatusOK)
}

func sendApplicationSection(w http.ResponseWriter, data types.ApplicationSection) error {
	return utils.SendJson(w, struct {
		Section types.ApplicationSection `json:"applicationSection"`
	}{Section: data}, http.StatusOK)
}
