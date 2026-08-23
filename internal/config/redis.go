package config

import "time"

type RedisConfig struct {
	Addr        string        `envconfig:"ADDR" required:"true"`
	Password    string        `envconfig:"PASSWORD" required:"true"`
	DB          int           `envconfig:"DB" default:"0"`
	DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s"`
}
