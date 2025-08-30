package spotify

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.mattglei.ch/musicsync/internal/apis"
	"go.mattglei.ch/musicsync/internal/secrets"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	ExpiresAt time.Time
}

func Authorize(client *http.Client) (AccessToken, error) {
	params := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {secrets.ENV.SpotifyClientID},
		"client_secret": {secrets.ENV.SpotifyClientSecret},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://accounts.spotify.com/api/token?%s", params.Encode()),
		nil,
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return AccessToken{}, fmt.Errorf("%w creating new request failed", err)
	}

	resp, err := apis.RequestJSON[AccessToken]("[spotify]", client, req)
	if err != nil {
		return AccessToken{}, fmt.Errorf("%w performing request failed", err)

	}

	resp.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn)*time.Second - 30*time.Second)

	return resp, nil
}
