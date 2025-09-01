package spotify

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mattglei.ch/musicsync/internal/apis"
	"go.mattglei.ch/timber"
)

func SendSpotifyAPIRequest[T any](
	client *http.Client,
	token *AccessToken,
	path string,
) (T, error) {
	var zeroValue T

	if token.ExpiresAt.Before(time.Now()) {
		newToken, err := Authorize(client)
		if err != nil {
			return zeroValue, fmt.Errorf("%w failed to refresh access token", err)
		}
		*token = newToken
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://api.spotify.com/%s", strings.TrimLeft(path, "/")),
		nil,
	)
	if err != nil {
		return zeroValue, fmt.Errorf("%w failed to create request", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	timber.Debug(token.Token)

	resp, err := apis.RequestJSON[T]("[spotify]", client, req)
	if err != nil {
		return zeroValue, fmt.Errorf("%w failed to make apple music API request", err)
	}
	return resp, nil
}
