package config

import "time"

type VKIDConfig struct {
	BaseURL         string        `envconfig:"BASE_URL"                 required:"true"`
	ClientID        string        `envconfig:"CLIENT_ID"                required:"true"`
	RedirectURL     string        `envconfig:"REDIRECT_URL"             required:"true"`
	Scope           string        `envconfig:"SCOPE"                    default:""`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"180d"`
}

type VKConfig struct {
	BaseURL string `envconfig:"BASE_URL" required:"true"`
}
