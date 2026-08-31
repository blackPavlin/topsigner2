package config

import "time"

type HttpConfig struct {
	Port              int           `envconfig:"PORT"                default:"8080"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT"        default:"5s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT"       default:"10s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT"        default:"120s"`
}

type CorsConfig struct {
	AllowedOrigins   []string `envconfig:"ALLOWED_ORIGINS"   required:"true"`
	AllowedMethods   []string `envconfig:"ALLOWED_METHODS"   required:"true"`
	AllowedHeaders   []string `envconfig:"ALLOWED_HEADERS"   required:"true"`
	ExposedHeaders   []string `envconfig:"EXPOSED_HEADERS"   required:"true"`
	AllowCredentials bool     `envconfig:"ALLOW_CREDENTIALS" required:"true"`
	MaxAge           int      `envconfig:"MAX_AGE"           required:"true"`
}
