package config

import (
	"fmt"
	"time"
)

type PostgresConfig struct {
	Host            string        `envconfig:"HOST"               required:"true"`
	Port            int           `envconfig:"PORT"               required:"true"`
	User            string        `envconfig:"USER"               required:"true"`
	Password        string        `envconfig:"PASSWORD"           required:"true"`
	Database        string        `envconfig:"DATABASE"           required:"true"`
	Schema          string        `envconfig:"SCHEMA"             required:"true"`
	SSLMode         string        `envconfig:"SSL_MODE"           required:"true"`
	MaxOpenConns    int32         `envconfig:"MAX_OPEN_CONNS"     default:"25"`
	MaxIdleConns    int32         `envconfig:"MAX_IDLE_CONNS"     default:"5"`
	ConnTimeout     time.Duration `envconfig:"CONN_TIMEOUT"       default:"5s"`
	ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME"  default:"1h"`
	ConnMaxIdleTime time.Duration `envconfig:"CONN_MAX_IDLE_TIME" default:"15m"`
}

func (c *PostgresConfig) ToDataSource() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s connect_timeout=%d",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode, c.Schema, int(c.ConnTimeout.Seconds()),
	)
}
