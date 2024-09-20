package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

var AccountIdKey = "account_id"

func Authenticator(s db.AuthStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			newRequest, err := authMiddlewareHandler(w, r, s)
			if err != nil {
				if error, ok := err.(types.ApiError); ok {
					utils.SendErrors(w, error.Errors, error.StatusCode)
				} else {
					log.Println("Internal server error: ", error.Error())
					utils.SendInternalServerError(w)
				}
			} else {
				next.ServeHTTP(w, newRequest)
			}
		}

		return http.HandlerFunc(hfn)
	}
}

func authMiddlewareHandler(w http.ResponseWriter, r *http.Request, s db.AuthStore) (*http.Request, error) {
	accountId, ok := verifyAccessToken(r)
	if ok {
		ctx := context.WithValue(r.Context(), AccountIdKey, accountId)
		r = r.WithContext(ctx)
		return r, nil
	}

	accountId, tokenVersion, ok := verifyRefreshToken(r, s)
	if !ok {
		utils.InvalidateAuthCookies(w)
		return nil, types.Unauthenticated
	}

	if err := utils.CreateAndSetAuthCookies(w, accountId, tokenVersion); err != nil {
		return nil, err
	}

	ctx := context.WithValue(r.Context(), AccountIdKey, accountId)
	r = r.WithContext(ctx)
	return r, nil
}

func verifyAccessToken(r *http.Request) (int, bool) {
	accessTokenCookie, err := r.Cookie(utils.AccessTokenName)
	if err != nil {
		return 0, false
	}
	return utils.VerifyAccessToken(accessTokenCookie.Value)
}

func verifyRefreshToken(r *http.Request, s db.AuthStore) (int, int, bool) {
	refreshTokenCookie, err := r.Cookie(utils.RefreshTokenName)
	if err != nil {
		return 0, 0, false
	}

	accountId, version, ok := utils.VerifyRefreshToken(refreshTokenCookie.Value)
	if !ok {
		return 0, 0, false
	}

	account, err := s.GetAccountById(accountId)
	if err != nil {
		return 0, 0, false
	}

	if account.TokenVersion != version {
		return 0, 0, false
	}

	return account.Id, account.TokenVersion, true
}
