package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store db.AuthStore
}

func NewHandler(store db.AuthStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) AddRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
		r.Post("/check", h.handleAuthCheck)
	})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body RegisterBody
	decoder.Decode(&body)

	if errors := body.IsValid(); len(errors) > 0 {
		service.SendErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	currentErrors := make([]string, 0)

	existingUser, err := h.store.GetUserByEmail(body.Email)
	if err != nil && !errors.Is(err, types.UserDoesNotExistErr) {
		service.SendInternalServerError(w)
		return
	}

	if existingUser.Email == body.Email {
		currentErrors = append(currentErrors, "User with provided email already exists")
	}

	if len(currentErrors) > 0 {
		service.SendErrorsResponse(w, currentErrors, http.StatusBadRequest)
		return
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	user, err := h.store.CreateUser(body.CreateUser(hash))
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	token, err := service.JwtClient.CreateToken(user.Id)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	type newUserResponseData struct {
		User  types.User `json:"user"`
		Token string     `json:"auth_token"`
	}

	data := newUserResponseData{
		User:  user,
		Token: token,
	}

	service.SendJsonResponse(w, data, http.StatusOK)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body LoginBody
	decoder.Decode(&body)

	if err := body.IsValid(); err != nil {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
		return
	}

	existingUser, err := h.store.GetUserByEmail(*body.Email)
	if errors.Is(err, types.UserDoesNotExistErr) {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusMethodNotAllowed)
		return
	} else if err != nil {
		service.SendInternalServerError(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(*body.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		service.SendErrorsResponse(w, []string{"Passwords do not match"}, http.StatusBadRequest)
		return
	} else if err != nil {
		service.SendInternalServerError(w)
		return
	}

	token, err := service.JwtClient.CreateToken(existingUser.Id)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	type loginResponseData struct {
		Token string `json:"auth_token"`
	}

	data := loginResponseData{
		Token: token,
	}

	service.SendJsonResponse(w, data, http.StatusOK)
}

func (h *Handler) handleAuthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := service.JwtClient.VerifyToken(r)
	if err != nil {
		service.SendErrorsResponse(w, []string{"Athorization header is not valid"}, http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
