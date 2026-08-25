package config

import "github.com/bboykiv/topsigner/internal/model"

type UserConfig struct {
	Default DefaultUserConfig `envconfig:"DEFAULT"`
}

type DefaultUserConfig struct {
	Email    string     `envconfig:"EMAIL"    default:"admin@topsigner.com"`
	Password string     `envconfig:"PASSWORD" default:"password123$"`
	Role     model.Role `envconfig:"ROLE"     default:"ADMIN"`
}
