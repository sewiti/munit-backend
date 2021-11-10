package config

import "time"

type Munit struct {
	Addr          string        `envconfig:"default=:7878"`
	AllowedOrigin string        `envconfig:"default=munit.digital"`
	Debug         bool          `envconfig:"default=false"`
	DSN           string        // Data source name
	SecretFile    string        `envconfig:"default=.secret"`
	Timeout       time.Duration `envconfig:"default=30s"`
}
