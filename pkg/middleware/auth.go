package middleware

import (
	"context"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
)

func JwtAuthenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			userId, err := service.JwtClient.VerifyToken(r)
			if err != nil || userId < 0 {
				service.SendErrorsResponse(w, []string{"Athorization header is not valid"}, http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, service.UserIdKey, userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
