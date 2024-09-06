package middleware

import (
	"net/http"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/go-chi/jwtauth/v5"
)

func JwtAuthenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("authorization")

			if token == "" {
				service.SendErrorsResponse(w, []string{"Athorization header not provided"}, http.StatusUnauthorized)
				return
			}

			_, err := jwtauth.VerifyToken(service.TokenAuth, strings.TrimPrefix(token, "Bearer "))
			if !strings.HasPrefix(token, "Bearer ") || err != nil {
				service.SendErrorsResponse(w, []string{"Athorization header is not valid"}, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
