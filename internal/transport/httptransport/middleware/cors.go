package middleware

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/bboykiv/topsigner/internal/config"
)

func Cors(config *config.Config) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		ExposedHeaders:   config.Cors.ExposedHeaders,
		AllowCredentials: config.Cors.AllowCredentials,
		MaxAge:           config.Cors.MaxAge,
	})
}
