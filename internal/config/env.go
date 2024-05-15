package config

import (
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
