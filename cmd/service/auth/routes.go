package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/cmd/db"
	"github.com/CelanMatjaz/job_application_tracker_api/cmd/service"
	"github.com/CelanMatjaz/job_application_tracker_api/cmd/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

var tokenAuth *jwtauth.JWTAuth

type Handler struct {
	store db.AuthStore
}

func NewHandler(store db.AuthStore) *Handler {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

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

	if err := body.IsValid(); err != nil {
		service.SendErrorsResponse(w, []string{err.Error()}, http.StatusBadRequest)
		return
	}

	currentErrors := make([]string, 0)

	existingUser, err := h.store.GetInternalUserByEmail(*body.Email)
	if err != nil && !errors.Is(err, UserDoesNotExistErr) {
		service.SendInternalServerError(w)
		return
	}

	if existingUser.Email == *body.Email {
		currentErrors = append(currentErrors, "User with provided email already exists")
	}

	if len(currentErrors) > 0 {
		service.SendErrorsResponse(w, currentErrors, http.StatusBadRequest)
		return
	}

	hash, err := hashPassword(*body.Password)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	user, err := h.store.CreateUser(body.CreateInternalUser(hash))
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	token, err := createJwtToken(existingUser)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	type newUserResponseData struct {
		User  types.User `json:"user"`
		Token string     `json:"auth_token"`
	}

	data := newUserResponseData{
		User:  user.User,
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

	existingUser, err := h.store.GetInternalUserByEmail(*body.Email)
	if errors.Is(err, UserDoesNotExistErr) {
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

	token, err := createJwtToken(existingUser)
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
	token := r.Header.Get("authorization")

	isValid := strings.HasPrefix(token, "Bearer ")
	_, err := jwtauth.VerifyToken(tokenAuth, strings.TrimPrefix(token, "Bearer "))
	if !isValid || token == "" || err != nil {
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

func createJwtToken(user types.InternalUser) (string, error) {
	_, token, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": user.Id,
	})
	return token, err
}
