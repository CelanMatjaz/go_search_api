package handlers

import (
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type ResumeHandler struct {
	store db.ResumeStore
}

func NewResumeHandler(store db.ResumeStore) *ResumeHandler {
	return &ResumeHandler{store: store}
}

func (h *ResumeHandler) AddRoutes(r chi.Router) {
	r.Route("/resumes", func(r chi.Router) {
		r.Route("/presets", func(r chi.Router) {
			r.Get("/", CreateHandler(CreateGenericGetManyWithPaginationHandler(h.store.GetResumePresets, sendResumePresets)))
			r.Get("/{id}", CreateHandler(CreateGenericGetSingleHandler(h.store.GetResumePreset, sendResumePreset)))
			r.Post("/", CreateHandler(CreateGenericPostHandler(h.store.CreateResumePreset, sendResumePreset)))
			r.Put("/{id}", CreateHandler(CreateGenericPutHandler(h.store.UpdateResumePreset, sendResumePreset)))
			r.Delete("/{id}", CreateHandler(CreateGenericDeleteHandler(h.store.DeleteResumePreset)))
		})

		r.Route("/sections", func(r chi.Router) {
			r.Get("/", CreateHandler(CreateGenericGetManyWithPaginationHandler(h.store.GetResumeSections, sendResumeSections)))
			r.Get("/{id}", CreateHandler(CreateGenericGetSingleHandler(h.store.GetResumeSection, sendResumeSection)))
			r.Post("/", CreateHandler(CreateGenericPostHandler(h.store.CreateResumeSection, sendResumeSection)))
			r.Put("/{id}", CreateHandler(CreateGenericPutHandler(h.store.UpdateResumeSection, sendResumeSection)))
			r.Delete("/{id}", CreateHandler(CreateGenericDeleteHandler(h.store.DeleteResumeSection)))
		})
	})
}

func sendResumePresets(w http.ResponseWriter, data []ResPreWithTags) error {
	return utils.SendJson(w, struct {
		Presets []ResPreWithTags `json:"resumePresets"`
	}{Presets: data}, http.StatusOK)
}

func sendResumePreset(w http.ResponseWriter, data types.ResumePreset) error {
	return utils.SendJson(w, struct {
		Preset types.ResumePreset `json:"resumePreset"`
	}{Preset: data}, http.StatusOK)
}

func sendResumeSections(w http.ResponseWriter, data []ResSecWithTags) error {
	return utils.SendJson(w, struct {
		Sections []ResSecWithTags `json:"resumeSections"`
	}{Sections: data}, http.StatusOK)
}

func sendResumeSection(w http.ResponseWriter, data types.ResumeSection) error {
	return utils.SendJson(w, struct {
		Section types.ResumeSection `json:"resumeSection"`
	}{Section: data}, http.StatusOK)
}
