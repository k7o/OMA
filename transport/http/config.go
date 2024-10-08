package http

import "fmt"

type Config struct {
	Port int `envconfig:"default=8080"`
}

func (c *Config) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("port is required")
	}

	return nil
}
