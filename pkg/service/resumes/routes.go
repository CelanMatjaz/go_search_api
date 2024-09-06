package resumes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
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
	r.Route("/resumes", func(r chi.Router) {
		r.Get("/", h.handleResumes)
		r.Get("/{resumeId}", h.handleSingleResume)
		r.Post("/", h.handlePostResume)
		r.Put("/{resumeId}", h.handlePutResume)
		r.Delete("/{resumeId}", h.handleDeleteResume)
	})
}

func (h *Handler) handleResumes(w http.ResponseWriter, r *http.Request) {
	pagination := service.GetPaginationParams(r)

	userId := r.Context().Value(service.UserIdKey)
	if userId == 0 {
		service.SendErrorsResponse(w, []string{"Unauthorized"}, http.StatusUnauthorized)
		return
	}

	resumes, err := h.store.GetRecords(userId, pagination.GetOffset(), pagination.Count)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, resumes, http.StatusOK)
}

func (h *Handler) handleSingleResume(w http.ResponseWriter, r *http.Request) {
	resumeId, err := strconv.Atoi(chi.URLParam(r, "resumeId"))
	if err != nil {
		service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(service.UserIdKey).(int)
	if userId == 0 {
		service.SendErrorsResponse(w, []string{"Unauthorized"}, http.StatusUnauthorized)
		return
	}

	resume, err := h.store.GetRecord(resumeId, userId)
	if errors.Is(err, sql.ErrNoRows) {
		service.SendErrorsResponse(w, []string{"Resume does not exist"}, http.StatusNotFound)
		return
	}
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, resume, http.StatusOK)
}

func (h *Handler) handlePostResume(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body ResumePostBody
	decoder.Decode(&body)

	if err := body.IsValid(); err != nil {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(service.UserIdKey).(int)
	if userId == 0 {
		service.SendErrorsResponse(w, []string{"Unauthorized"}, http.StatusUnauthorized)
		return
	}

	newResume, err := h.store.CreateRecord(userId, *body.Name, *body.Note)
	if err != nil {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
		return
	}

	service.SendJsonResponse(w, newResume, http.StatusOK)
}

func (h *Handler) handlePutResume(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body ResumePostBody
	decoder.Decode(&body)

	resumeId, err := strconv.Atoi(chi.URLParam(r, "resumeId"))
	if err != nil {
		service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
		return
	}

	if err := body.IsValid(); err != nil {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(service.UserIdKey).(int)
	if userId == 0 {
		service.SendErrorsResponse(w, []string{"Unauthorized"}, http.StatusUnauthorized)
		return
	}

	newResume, err := h.store.UpdateRecord(resumeId, userId, *body.Name, *body.Note)
	if errors.Is(err, sql.ErrNoRows) {
		service.SendErrorsResponse(w, []string{"Resume with provided id was not found"}, http.StatusNotFound)
		return
	}
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, newResume, http.StatusOK)
}

func (h *Handler) handleDeleteResume(w http.ResponseWriter, r *http.Request) {
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
