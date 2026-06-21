package config

type HttpConfig struct {
	Port int `envconfig:"PORT" default:"8080"`
}

type CorsConfig struct {
	AllowedOrigins   []string `envconfig:"ALLOWED_ORIGINS" required:"true"`
	AllowedMethods   []string `envconfig:"ALLOWED_METHODS" required:"true"`
	AllowedHeaders   []string `envconfig:"ALLOWED_HEADERS" required:"true"`
	ExposedHeaders   []string `envconfig:"EXPOSED_HEADERS" required:"true"`
	AllowCredentials bool     `envconfig:"ALLOW_CREDENTIALS" required:"true"`
	MaxAge           int      `envconfig:"MAX_AGE" required:"true"`
}
