package applications

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
	store db.ApplicationStore
}

func NewHandler(store db.ApplicationStore) *Handler {
	return &Handler{store: store}
}

type ApplicationsResponse struct {
	Data []types.Application `json:"applications"`
}

type ApplicationResponse struct {
	Data types.Application `json:"application"`
}

func (h *Handler) AddRoutes(r chi.Router) {
	r.Route("/applications", func(r chi.Router) {
		r.Use(middleware.JwtAuthenticator())

		r.Get("/", h.handleGetApplications)
		r.Get("/{id}", h.handleGetApplication)
		r.Post("/", h.handlePostApplication)
		r.Put("/{id}", h.handlePutApplication)
		r.Delete("/{id}", h.handleDeleteApplication)
	})
}

func (h *Handler) handleGetApplications(w http.ResponseWriter, r *http.Request) {
	pagination := service.GetPaginationParams(r)
	userId := r.Context().Value(service.UserIdKey).(int)

	applications, err := h.store.GetUserApplications(userId, pagination)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, ApplicationsResponse{applications}, http.StatusOK)
}

func (h *Handler) handleGetApplication(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(service.UserIdKey).(int)
	applicationId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
		return
	}

	application, err := h.store.GetUserApplication(userId, applicationId)
    println(err.Error())
	switch err {
	case nil:
		service.SendJsonResponse(w, ApplicationResponse{application}, http.StatusOK)
		return
	case types.RecordDoesNotExist:
		service.SendErrorsResponse(w, []string{"Application does not exist"}, http.StatusNotFound)
		return
	default:
		service.SendInternalServerError(w)
		return
	}
}

func (h *Handler) handlePostApplication(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body types.Application
	decoder.Decode(&body)

	if errors := body.IsValid(); len(errors) > 0 {
		service.SendErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(service.UserIdKey).(int)
	newApplication, err := h.store.CreateUserApplication(userId, body)
	switch err {
	case nil:
		service.SendJsonResponse(w, ApplicationResponse{newApplication}, http.StatusOK)
		return
	default:
		service.SendInternalServerError(w)
		return
	}
}

func (h *Handler) handlePutApplication(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body types.Application
	decoder.Decode(&body)

	if errors := body.IsValid(); len(errors) > 0 {
		service.SendErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(service.UserIdKey).(int)
	newApplication, err := h.store.UpdateUserApplication(userId, body)
	switch err {
	case nil:
		service.SendJsonResponse(w, ApplicationResponse{newApplication}, http.StatusOK)
		return
	default:
		service.SendInternalServerError(w)
		return
	}
}

func (h *Handler) handleDeleteApplication(w http.ResponseWriter, r *http.Request) {
	applicationId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		service.SendErrorsResponse(w, []string{"Provided url param was not parsable"}, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(service.UserIdKey).(int)

	err = h.store.DeleteUserApplication(userId, applicationId)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, nil, http.StatusOK)
}
