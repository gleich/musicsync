package spotify

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.mattglei.ch/musicsync/internal/apis"
	"go.mattglei.ch/musicsync/internal/secrets"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    time.Time
}

func (c *SpotifyClient) Authorize() error {
	params := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.Tokens.RefreshToken},
		"client_id":     {secrets.ENV.SpotifyClientID},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://accounts.spotify.com/api/token?%s", params.Encode()),
		nil,
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(
		secrets.ENV.SpotifyClientID+":"+secrets.ENV.SpotifyClientSecret,
	)))
	if err != nil {
		return fmt.Errorf("%w creating new request failed", err)
	}

	resp, err := apis.RequestJSON[Tokens]("[spotify]", c.HttpClient, req)
	if err != nil {
		return fmt.Errorf("%w performing request failed", err)

	}

	resp.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn)*time.Second - 30*time.Second)

	c.mutex.Lock()
	c.Tokens = &resp
	c.mutex.Unlock()
	return nil
}
