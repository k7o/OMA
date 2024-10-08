package config

import (
	"fmt"
	"oma/internal/revision"
	"oma/transport/http"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel       zerolog.Level           `envconfig:"default=1"`
	RevisionConfig revision.RevisionConfig `envconfig:"optional"`
	OpaDownloadUrl string
	Transport      TransportConfig
}

type TransportConfig struct {
	HTTP http.Config
}

func (c *Config) Validate() error {
	if err := c.RevisionConfig.Validate(); err != nil {
		return err
	}

	if c.OpaDownloadUrl == "" {
		return fmt.Errorf("opa_download_url is required")
	}

	if err := c.Transport.HTTP.Validate(); err != nil {
		return err
	}

	return nil
}
