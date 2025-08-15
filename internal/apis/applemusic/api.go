package applemusic

import (
	"fmt"
	"net/http"
	"strings"

	"go.mattglei.ch/musicsync/internal/apis"
	"go.mattglei.ch/musicsync/internal/secrets"
)

func SendAppleMusicAPIRequest[T any](client *http.Client, path string) (T, error) {
	var zeroValue T
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://api.music.apple.com/%s", strings.TrimLeft(path, "/")),
		nil,
	)
	if err != nil {
		return zeroValue, fmt.Errorf("%w failed to create request", err)
	}
	req.Header.Set("Authorization", "Bearer "+secrets.ENV.AppleMusicAppToken)
	req.Header.Set("Music-User-Token", secrets.ENV.AppleMusicUserToken)

	resp, err := apis.RequestJSON[T]("[apple music]", client, req)
	if err != nil {
		return zeroValue, fmt.Errorf("%w failed to make apple music API request", err)
	}
	return resp, nil
}
