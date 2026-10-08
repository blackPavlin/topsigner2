package config

import "github.com/bboykiv/topsigner/internal/model"

type UserConfig struct {
	Default DefaultUserConfig `envconfig:"DEFAULT"`
}

type DefaultUserConfig struct {
	Email    string     `envconfig:"EMAIL"    required:"true"`
	Password string     `envconfig:"PASSWORD" required:"true"`
	Role     model.Role `envconfig:"ROLE"     required:"true"`
}
