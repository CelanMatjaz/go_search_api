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
	cookie.Expires = time.Time{}
	cookie.MaxAge = -1
    cookie.Secure = false
	http.SetCookie(w, &cookie)
}

func InvalidateAuthCookies(w http.ResponseWriter) {
	InvalidateCookie(w, AccessTokenName)
	InvalidateCookie(w, RefreshTokenName)
}

func CreateAndSetAuthCookies(w http.ResponseWriter, accountId int, tokenVersion int) error {
	token, err := CreateAccessToken(accountId)
	if err != nil {
		return err
	}

	refreshToken, err := CreateRefreshToken(accountId, tokenVersion)
	if err != nil {
		return err
	}

	CreateAndSetCookie(w, AccessTokenName, token)
	CreateAndSetCookie(w, RefreshTokenName, refreshToken)

	return nil
}
