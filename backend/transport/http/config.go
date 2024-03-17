package http

type Config struct {
	Port int `envconfig:"PORT" default:"8080"`
}
