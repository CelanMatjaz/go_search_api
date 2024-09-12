package utils

import (
	"net/http"
	"time"
)

const cookieMaxAge = 60 * 60 * 24 * 30

func CreateCookie(name string, value string) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   cookieMaxAge,
		Path:     "/",
		// Domain: "go-search.site",
	}
}

func CreateAndSetCookie(w http.ResponseWriter, name string, value string) {
	cookie := CreateCookie(name, value)
	http.SetCookie(w, &cookie)
}

func InvalidateCookie(w http.ResponseWriter, name string) {
	cookie := CreateCookie(name, "")
	cookie.Expires = time.Now()
	http.SetCookie(w, &cookie)
}
