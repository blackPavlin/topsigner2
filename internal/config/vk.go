package config

import "time"

type VKIDConfig struct {
	BaseURL         string        `envconfig:"BASE_URL"          required:"true"`
	ClientID        string        `envconfig:"CLIENT_ID"         required:"true"`
	RedirectURL     string        `envconfig:"REDIRECT_URL"      required:"true"`
	SecretKey       string        `envconfig:"SECRET_KEY"        required:"true"`
	ServiceKey      string        `envconfig:"SERVICE_KEY"       required:"true"`
	Scope           string        `envconfig:"SCOPE"             default:""`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TOKEN_TTL" default:"4320h"`
}

type VKConfig struct {
	BaseURL          string `envconfig:"BASE_URL"           required:"true"`
	OAuthBaseURL     string `envconfig:"OAUTH_BASE_URL"     required:"true"`
	OAuthRedirectURL string `envconfig:"OAUTH_REDIRECT_URL" required:"true"`
	OAuthScope       string `envconfig:"OAUTH_SCOPE"        default:"manage,messages"`
}
