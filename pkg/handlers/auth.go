package handlers

import (
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	store db.AuthStore
}

func NewAuthHandler(store db.AuthStore) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) AddRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", CreateHandler(h.HandleRegister))
		r.Post("/login", CreateHandler(h.HandleLogin))
		r.Post("/logout", CreateHandler(h.HandleLogout))
		r.Post("/oauth", CreateHandler(h.HandleOAuth))

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticator(h.store))
			r.Post("/check", CreateHandler(h.HandleAuthCheck))
		})
	})
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	body, err := decodeAndValidateBody[types.RegisterBody](r)
	if err != nil {
		return err
	}

	if body.Password != body.PasswordVerify {
		return types.PasswordsDoNotMatch
	}

	_, exists, err := h.store.GetAccountByEmail(body.Email)
	if err != nil {
		return err
	}
	if exists {
		return types.AccountAlreadyExists
	}

	newAccountData, err := types.CreateNewAccountData(body.DisplayName, body.Email, body.Password)
	if err != nil {
		return err
	}

	_, err = h.store.CreateAccount(newAccountData)
	if err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	body, err := decodeAndValidateBody[types.LoginBody](r)
	if err != nil {
		return err
	}

	existingAccount, exists, err := h.store.GetAccountByEmail(body.Email)
	if err != nil {
		return err
	}
	if !exists {
		return types.AccountDoesNotExist
	}

	ok, err := body.ComparePassword(existingAccount.PasswordHash.String)
	if err != nil {
		return err
	}
	if !ok {
		return types.InvalidPassword
	}

	err = utils.CreateAndSetAuthCookies(w, existingAccount.Id, existingAccount.TokenVersion)
	if err != nil {
		return utils.SendInternalServerError(w)
	}

	return utils.SendJson(w, struct {
		Account types.Account `json:"accountData"`
	}{Account: existingAccount}, http.StatusOK)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, _ *http.Request) error {
	utils.InvalidateAuthCookies(w)
	return nil
}

func (h *AuthHandler) HandleAuthCheck(w http.ResponseWriter, r *http.Request) error {
	accountId, ok := r.Context().Value(middleware.AccountIdKey).(int)
	if !ok {
		return types.Unauthenticated
	}

	existingAccount, exists, err := h.store.GetAccountById(accountId)
	if err != nil {
		return err
	}
	if !exists {
		return types.AccountDoesNotExist
	}

	return utils.SendJson(w, struct {
		Account types.Account `json:"accountData"`
	}{Account: existingAccount}, http.StatusOK)
}

func (h *AuthHandler) HandleOAuth(w http.ResponseWriter, r *http.Request) error {
	body, err := decodeAndValidateBody[types.OAuthLoginBody](r)
	if err != nil {
		return types.UnparsableJsonBody
	}

	oauthClient, exists, err := h.store.GetOAuthClientByName(body.Provider)
	if err != nil {
		return err
	}
	if !exists {
		return types.UnknownOAuthProvider
	}

	tokenResponse, statusCode, err := utils.FetchToken(oauthClient, body.Code)
	if err != nil {
		return err
	} else if statusCode != 200 {
		return types.OAuthProviderIssues
	}

	infoResponse, statusCode, err := utils.FetchAccountData(oauthClient.AccountDataEndpoint, tokenResponse.AccessToken)
	if err != nil {
		return err
	} else if statusCode != 200 {
		return types.OAuthProviderIssues
	} else if !infoResponse.EmailVerified {
		return types.UnverifiedOAuthEmail
	}

	account, exists, err := h.store.GetAccountByEmail(infoResponse.Email)
	if err != nil {
		return err
	}

	if !exists {
		account, err = h.store.CreateAccountWithOAuth(types.Account{
			DisplayName: infoResponse.Name,
			Email:       infoResponse.Email,
		}, tokenResponse, oauthClient.Id)
	}

	err = utils.CreateAndSetAuthCookies(w, account.Id, account.TokenVersion)
	if err != nil {
		return utils.SendInternalServerError(w)
	}

	return utils.SendJson(w, struct {
		Account types.Account `json:"accountData"`
	}{Account: account}, http.StatusOK)
}
