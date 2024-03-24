package config

import (
	"oma/internal/revision"
	"oma/transport/http"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel       zerolog.Level
	RevisionConfig revision.RevisionConfig
	Transport      TransportConfig
}

type TransportConfig struct {
	HTTP http.Config
}
