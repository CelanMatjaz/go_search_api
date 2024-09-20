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
	Common
	Name                string `json:"-" db:"name"`
	ClientId            string `json:"-" db:"clientId"`
	ClientSecret        string `json:"-" db:"clientSecret"`
	Scopes              string `json:"-" db:"scopes"`
	CodeEndpoint        string `json:"-" db:"codeEndpoint"`
	TokenEndpoint       string `json:"-" db:"tokenEndpoint"`
	AccountDataEndpoint string `json:"-" db:"accountDataEndpoint"`
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
