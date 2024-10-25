package types

type OAuthLoginBody struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

func (b OAuthLoginBody) Verify() []string {
	if b.Code == "" || b.Provider == "" {
		return InvalidJsonBodyErrors
	}

	return nil
}

type OAuthClient struct {
	WithId
	Name                string `json:"-" db:"name"`
	ClientId            string `json:"-" db:"client_id"`
	ClientSecret        string `json:"-" db:"client_secret"`
	Scopes              string `json:"-" db:"scopes"`
	CodeEndpoint        string `json:"-" db:"code_endpoint"`
	TokenEndpoint       string `json:"-" db:"token_endpoint"`
	AccountDataEndpoint string `json:"-" db:"account_data_endpoint"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type GoogleAccountInfo struct {
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
