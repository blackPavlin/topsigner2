package config

import "time"

type AuthConfig struct {
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL"  required:"true"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL"  required:"true"`
	SigningKey      string        `envconfig:"SIGNING_KEY" required:"true"`
}
