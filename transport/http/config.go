package http

type Config struct {
	Port int `envconfig:"default=8080"`
}
