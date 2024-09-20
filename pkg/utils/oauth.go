package utils

import (
	"encoding/json"
	"strings"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func FetchToken(oauthClient types.OAuthClient, code string) (types.TokenResponse, int, error) {
	tokenJsonBody, _ := json.Marshal(map[string]string{
		"client_id":     oauthClient.ClientId,
		"client_secret": oauthClient.ClientSecret,
		"code":          code,
		"grant_type":    "authorization_code",
		// "redirect_uri":  "http://dev.go-search.site/callback",
		"redirect_uri": "http://localhost:3000/callback",
	})

	tokenResponse, statusCode, err := MakeRequest[types.TokenResponse](
		oauthClient.TokenEndpoint,
		"POST",
		map[string]string{"Content-Type": "application/json"},
		tokenJsonBody,
	)

	if err != nil {
		return tokenResponse, 0, err
	}

	return tokenResponse, statusCode, nil
}

func FetchAccountData(url string, token string) (types.GoogleAccountInfo, int, error) {
	infoResponse, statusCode, err := MakeRequest[types.GoogleAccountInfo](
		url, "GET",
		map[string]string{
			"Authorization": strings.Join([]string{"Bearer ", token}, ""),
		},
		[]byte{},
	)

	if err != nil {
		return infoResponse, 0, err
	}

	return infoResponse, statusCode, nil
}
