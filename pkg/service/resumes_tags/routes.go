package resumestags

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store db.GenericStore[types.Resume]
}

func NewHandler(store db.GenericStore[types.Resume]) *Handler {
	return &Handler{store: store}
}

func (h *Handler) AddRoutes(r chi.Router) {
	r.Route("/resume_tags", func(r chi.Router) {
		r.Use(middleware.JwtAuthenticator(service.TokenAuth))
		r.Get("/{resumeId}", h.handleTags)
		r.Post("/", h.handlePostTag)
		r.Put("/{resumeId}", h.handlePutTag)
		r.Delete("/{resumeId}", h.handleDeleteTag)
	})
}

func (h *Handler) handleTags(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(chi.URLParam(r, "resumeId"))
	if err != nil {
		service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
		return
	}

	tags, err := h.store.GetRecords()
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, tags, http.StatusOK)
}

func (h *Handler) handlePostTag(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body ResumeTagPostBody
	decoder.Decode(&body)

	if err := body.IsValid(); err != nil {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
		return
	}

	// newResume, err := h.store.CreateRecord(1, *body.Name, *body.Note)
	// if err != nil {
	// 	service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
	// 	return
	// }
	//
	// service.SendJsonResponse(w, newResume, http.StatusOK)
}

func (h *Handler) handlePutTag(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body ResumeTagPostBody
	decoder.Decode(&body)

	// resumeId, err := strconv.Atoi(chi.URLParam(r, "resumeId"))
	// if err != nil {
	// 	service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
	// 	return
	// }
	//
	// if err := body.IsValid(); err != nil {
	// 	service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
	// 	return
	// }

	// newResume, err := h.store.UpdateRecord(resumeId, *body.Name, *body.Note)
	// if errors.Is(err, sql.ErrNoRows) {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// if err != nil {
	// 	service.SendInternalServerError(w)
	// 	return
	// }
	//
	// service.SendJsonResponse(w, newResume, http.StatusOK)
}

func (h *Handler) handleDeleteTag(w http.ResponseWriter, r *http.Request) {
	resumeId, err := strconv.Atoi(chi.URLParam(r, "resumeId"))
	if err != nil {
		service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
		return
	}

	err = h.store.DeleteRecord(resumeId)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
