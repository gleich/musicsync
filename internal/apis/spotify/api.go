package spotify

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mattglei.ch/musicsync/internal/apis"
)

type SpotifyClient struct {
	HttpClient *http.Client
	Tokens     *Tokens
	mutex      sync.RWMutex
}

type spotifyRequest struct {
	Method string
	Path   string
	Body   io.Reader
}

func sendSpotifyAPIRequest[T any](
	client *SpotifyClient,
	request spotifyRequest,
) (T, error) {
	var zeroValue T

	if client.Tokens.ExpiresAt.Before(time.Now()) {
		err := client.Authorize()
		if err != nil {
			return zeroValue, fmt.Errorf("%w failed to refresh access token", err)
		}
	}

	req, err := http.NewRequest(
		request.Method,
		fmt.Sprintf("https://api.spotify.com/%s", strings.TrimLeft(request.Path, "/")),
		request.Body,
	)
	if err != nil {
		return zeroValue, fmt.Errorf("%w failed to create request", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Tokens.AccessToken))

	resp, err := apis.RequestJSON[T]("[spotify]", client.HttpClient, req)
	if err != nil {
		return zeroValue, fmt.Errorf("%w failed to make apple music API request", err)
	}
	return resp, nil
}
