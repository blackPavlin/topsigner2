package config

import "time"

type AuthConfig struct {
	ExpiresIn  time.Duration `envconfig:"EXPIRES_IN"  required:"true"`
	SigningKey string        `envconfig:"SIGNING_KEY" required:"true"`
}
