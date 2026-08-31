package config

import "time"

type AuthConfig struct {
	AccessTokenTTL  time.Duration `envconfig:"ACCESS_TOKEN_TTL"  required:"true"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" required:"true"`
	SigningKey      string        `envconfig:"SIGNING_KEY"       required:"true"`
}

type VKIDConfig struct {
	BaseURL     string `envconfig:"BASE_URL"     required:"true"`
	ClientID    string `envconfig:"CLIENT_ID"    required:"true"`
	RedirectURL string `envconfig:"REDIRECT_URL" required:"true"`
	Scope       string `envconfig:"SCOPE"        required:"true"`
}
