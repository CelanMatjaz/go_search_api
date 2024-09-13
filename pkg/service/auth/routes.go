package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

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
		r.Post("/logout", h.handleLogout)
		r.Post("/oauth", h.handleOAuth)

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
	providerDoesNotExist       = []string{"Provider does not exist"}
	errorAuthenticating        = []string{"Error authenticating user account with OAuth"}
	emailNotVerified           = []string{"Email not verified"}
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

	if err != nil {
		switch err {
		case types.UserDoesNotExistErr:
			break
		default:
			service.SendInternalServerError(w)
			return
		}
	}

	if existingUser.Email == body.Email {
		service.SendErrorsResponse(w, userWithEmailAlreadyExists, http.StatusBadRequest)
        return
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

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash.String), []byte(*body.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		service.SendErrorsResponse(w, passwordsDoNotMatch, http.StatusBadRequest)
		return
	} else if err != nil {
		service.SendInternalServerError(w)
		return
	}

	err = createAndSetCookies(w, existingUser)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, struct {
		User types.User `json:"user"`
	}{User: existingUser}, http.StatusOK)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	utils.InvalidateCookie(w, utils.AccessTokenName)
	utils.InvalidateCookie(w, utils.RefreshTokenName)
	service.SendJsonResponse(w, nil, http.StatusOK)
}

func (h *Handler) handleAuthCheck(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(utils.UserIdKey).(int)
	user, err := h.store.GetUserById(userId)

	if err != nil {
		service.SendErrorsResponse(w, userNotAuthenticated, http.StatusUnauthorized)
		return
	}

	service.SendJsonResponse(w, struct {
		User types.User `json:"user"`
	}{User: user}, http.StatusOK)
}

func (h *Handler) handleOAuth(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body struct {
		Code     string `json:"code"`
		Provider string `json:"provider"`
	}
	decoder.Decode(&body)

	oauthClient, err := h.store.GetOauthClientByName(body.Provider)
	if err != nil {
		switch err {
		case types.RecordDoesNotExist:
			service.SendErrorsResponse(w, providerDoesNotExist, http.StatusNotFound)
			return
		default:
			service.SendInternalServerError(w)
			return
		}
	}

	tokenResponse, statusCode, err := fetchToken(oauthClient, body.Code)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}
	if statusCode != 200 {
		service.SendErrorsResponse(w, errorAuthenticating, http.StatusUnauthorized)
		return
	}

	infoResponse, statusCode, err := fetchUserData(oauthClient.UserDataEndpoint, tokenResponse.AccessToken)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}
	if statusCode != 200 {
		service.SendErrorsResponse(w, errorAuthenticating, http.StatusUnauthorized)
		return
	}
	if infoResponse.EmailVerified == false {
		service.SendErrorsResponse(w, emailNotVerified, http.StatusUnauthorized)
		return
	}

	user, err := h.store.GetUserByEmail(infoResponse.Email)
	if err != nil {
		switch err {
		case types.UserDoesNotExistErr:
			user, err = h.store.CreateUserWithOAuth(types.User{
				DisplayName: infoResponse.Name,
				Email:       infoResponse.Email,
			}, tokenResponse, oauthClient.Id)

			if err == nil {
				break
			}
		default:
			service.SendInternalServerError(w)
			return
		}
	}

	err = createAndSetCookies(w, user)
	if err != nil {
		service.SendInternalServerError(w)
		return
	}

	service.SendJsonResponse(w, struct {
		User types.User `json:"user"`
	}{User: user}, http.StatusOK)
}

func makeRequest[ResponseBody any](url string, method string, headers map[string]string, requestBody []byte) (ResponseBody, int, error) {
	var body ResponseBody
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return body, 0, err
	}

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return body, 0, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return body, 0, err
	}

	err = json.Unmarshal(responseBody, &body)
	return body, response.StatusCode, err
}

func fetchToken(oauthClient types.OAuthClient, code string) (types.TokenResponse, int, error) {
	tokenJsonBody, _ := json.Marshal(map[string]string{
		"client_id":     oauthClient.ClientId,
		"client_secret": oauthClient.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		// "redirect_uri":  "http://dev.go-search.site/callback",
		"redirect_uri": "http://localhost:3000/callback",
	})

	tokenResponse, statusCode, err := makeRequest[types.TokenResponse](
		oauthClient.TokenEndpoint,
		"POST",
		map[string]string{"Content-Type": "application/json"},
		tokenJsonBody,
	)

	if err != nil {
		return tokenResponse, 0, err
	}

	return tokenResponse, statusCode, nil
}

func fetchUserData(url string, token string) (types.GoogleUserInfo, int, error) {
	infoResponse, statusCode, err := makeRequest[types.GoogleUserInfo](
		url, "GET",
		map[string]string{
			"Authorization": strings.Join([]string{"Bearer ", token}, ""),
		},
		[]byte{},
	)

	if err != nil {
		return infoResponse, 0, err
	}

	return infoResponse, statusCode, nil
}

func createAndSetCookies(w http.ResponseWriter, user types.User) error {
	token, err := utils.JwtClient.CreateAccessToken(user.Id)
	if err != nil {
		return err
	}

	refreshToken, err := utils.JwtClient.CreateRefreshToken(user.Id, user.TokenVersion)
	if err != nil {
		return err
	}

	utils.CreateAndSetCookie(w, utils.AccessTokenName, token)
	utils.CreateAndSetCookie(w, utils.RefreshTokenName, refreshToken)

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
