package secrets

import (
	"errors"
	"io/fs"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"go.mattglei.ch/timber"
)

var ENV Secrets

type Secrets struct {
	AppleMusicAppToken  string `env:"APPLE_MUSIC_APP_TOKEN"`
	AppleMusicUserToken string `env:"APPLE_MUSIC_USER_TOKEN"`
}

func Load() {
	if _, err := os.Stat(".env"); !errors.Is(err, fs.ErrNotExist) {
		err := godotenv.Load()
		if err != nil {
			timber.Fatal(err, "loading .env file failed")
		}
	}

	secrets, err := env.ParseAs[Secrets]()
	if err != nil {
		timber.Fatal(err, "parsing required env vars failed")
	}
	ENV = secrets
	timber.Done("loaded secrets")
}
