package config

import (
	"oma/transport/http"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel  zerolog.Level
	Transport TransportConfig
}

type TransportConfig struct {
	HTTP http.Config
}
