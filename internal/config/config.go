package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Auth     AuthConfig     `envconfig:"AUTH"`
	VKID     VKIDConfig     `envconfig:"VKID"`
	Http     HttpConfig     `envconfig:"HTTP"`
	Cors     CorsConfig     `envconfig:"CORS"`
	S3       S3Config       `envconfig:"S3"`
	Postgres PostgresConfig `envconfig:"POSTGRES"`
	Debug    bool           `envconfig:"DEBUG" default:"false"`
}

func New() (*Config, error) {
	config := new(Config)

	if err := envconfig.Process("", config); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return config, nil
}
