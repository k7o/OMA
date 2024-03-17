package config

import "oma/transport/http"

type Config struct {
	Transport TransportConfig
}

type TransportConfig struct {
	HTTP http.Config
}
