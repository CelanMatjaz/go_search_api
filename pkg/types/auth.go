package types

import "database/sql"

type User struct {
	Common
	DisplayName  string         `json:"display_name" db:"display_name"`
	Email        string         `json:"email" db:"email"`
	PasswordHash sql.NullString `json:"-" db:"password_hash"`
	TokenVersion int            `json:"-" db:"refresh_token_version"`
	Timestamps
}

type OAuthClient struct {
	Common
	Name             string `json:"-" db:"name"`
	ClientId         string `json:"-" db:"client_id"`
	ClientSecret     string `json:"-" db:"client_secret"`
	Scopes           string `json:"-" db:"scopes"`
	CodeEndpoint     string `json:"-" db:"code_endpoint"`
	TokenEndpoint    string `json:"-" db:"token_endpoint"`
	UserDataEndpoint string `json:"-" db:"user_data_endpoint"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type GoogleUserInfo struct {
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
