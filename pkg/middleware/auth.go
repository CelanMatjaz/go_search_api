package middleware

import (
	"context"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service/auth"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/utils"
)

var (
	invalidToken    = []string{"Invalid access token"}
	missingCookie   = []string{"Missing cookie"}
	unauthenticated = []string{"Unauthenticated"}
)

func Authenticator(s *auth.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			userId, err := verifyAccessToken(r)
			if err == nil {
				ctx := r.Context()
				ctx = context.WithValue(ctx, utils.UserIdKey, userId)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if handled := handleAccessTokenVerificationError(w, r, err); handled {
				return
			}

			userId, version, err := verifyRefreshToken(r, s)
			if handled := handleRefreshTokenVerificationError(w, r, err); handled {
				return
			}

			if handled := setCookies(w, userId, version); handled {
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, utils.UserIdKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func handleAccessTokenVerificationError(w http.ResponseWriter, r *http.Request, err error) bool {
	switch err {
	case nil:
	case types.InvalidTokenErr:
	case types.UnauthenticatedErr:
		return false
	case types.MissingCookieErr:
		service.SendErrorsResponse(w, missingCookie, http.StatusUnauthorized)
		return true
	default:
		service.SendInternalServerError(w)
		return true
	}

	return false
}

func handleRefreshTokenVerificationError(w http.ResponseWriter, r *http.Request, err error) bool {
	switch err {
	case nil:
	case types.InvalidTokenErr:
		return false
	case types.MissingCookieErr:
		service.SendErrorsResponse(w, invalidToken, http.StatusUnauthorized)
		return true
	case types.UnauthenticatedErr:
		service.SendErrorsResponse(w, unauthenticated, http.StatusUnauthorized)
		return true
	default:
		service.SendInternalServerError(w)
		return true
	}

	return false
}

func setCookies(w http.ResponseWriter, userId int, version int) bool {
	newAccessToken, err := utils.JwtClient.CreateAccessToken(userId)
	if err != nil {
		service.SendInternalServerError(w)
		return true
	}
	utils.CreateAndSetCookie(w, utils.AccessTokenName, newAccessToken)

	newRefreshToken, _ := utils.JwtClient.CreateRefreshToken(userId, version)
	if err != nil {
		service.SendInternalServerError(w)
		return true
	}
	utils.CreateAndSetCookie(w, utils.RefreshTokenName, newRefreshToken)

	return false
}

func verifyAccessToken(r *http.Request) (int, error) {
	accessTokenCookie, err := r.Cookie(utils.AccessTokenName)
	if err != nil {
		return 0, types.MissingCookieErr
	}
	return utils.JwtClient.VerifyAccessToken(accessTokenCookie.Value)
}

func verifyRefreshToken(r *http.Request, s *auth.Store) (int, int, error) {
	refreshTokenCookie, err := r.Cookie(utils.RefreshTokenName)
	if err != nil {
		return 0, 0, types.MissingCookieErr
	}
	userId, version, err := utils.JwtClient.VerifyRefreshToken(refreshTokenCookie.Value)
	if err != nil {
		return 0, 0, types.UnauthenticatedErr
	}

	user, err := s.GetUserById(userId)
	if err != nil {
		return 0, 0, err
	}

	if user.TokenVersion != version {
		return 0, 0, types.InvalidTokenErr
	}

	return user.Id, user.TokenVersion, nil
}
