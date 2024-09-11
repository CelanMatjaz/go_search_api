package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store db.AuthStore
}

func NewHandler(store db.AuthStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) AddRoutes(r chi.Router, auth func(http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)

		r.Group(func(r chi.Router) {
			r.Use(auth)
			r.Post("/check", h.handleAuthCheck)
		})
	})
}

var (
	userWithEmailAlreadyExists = []string{"User with provided email already exists"}
	userWithEmailDoesNotExist  = []string{"User with provided email does not exist"}
	passwordsDoNotMatch        = []string{"Passwords do not match"}
	userNotAuthenticated       = []string{"User not authenticated"}
)

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body RegisterBody
	decoder.Decode(&body)

	if errors := body.IsValid(); len(errors) > 0 {
		service.SendErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	existingUser, err := h.store.GetUserByEmail(body.Email)

	switch err {
	case types.UserDoesNotExistErr:
	case nil:
		break
	default:
		service.SendInternalServerError(w)
		return
	}

	if existingUser.Email == body.Email {
		service.SendErrorsResponse(w, userWithEmailAlreadyExists, http.StatusBadRequest)
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	_, err = h.store.CreateUser(body.CreateUser(hash))
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, nil, http.StatusOK)
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
		service.SendErrorsResponse(w, userWithEmailDoesNotExist, http.StatusUnauthorized)
		return
	} else if err != nil {
		service.SendInternalServerError(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(*body.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		service.SendErrorsResponse(w, passwordsDoNotMatch, http.StatusBadRequest)
		return
	} else if err != nil {
		service.SendInternalServerError(w)
		return
	}

	token, err := utils.JwtClient.CreateAccessToken(
		existingUser.Id,
	)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	refreshToken, err := utils.JwtClient.CreateRefreshToken(
		existingUser.Id,
		existingUser.TokenVersion,
	)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	utils.CreateAndSetCookie(w, utils.AccessTokenName, token)
	utils.CreateAndSetCookie(w, utils.RefreshTokenName, refreshToken)

	service.SendJsonResponse(w, struct {
		User types.User `json:"user"`
	}{User: existingUser}, http.StatusOK)
}

func (h *Handler) handleAuthCheck(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(utils.UserIdKey).(int)
	println(userId)
	user, err := h.store.GetUserById(userId)

	if err != nil {
		service.SendErrorsResponse(w, userNotAuthenticated, http.StatusUnauthorized)
		return
	}

	service.SendJsonResponse(w, struct {
		User types.User `json:"user"`
	}{User: user}, http.StatusOK)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
