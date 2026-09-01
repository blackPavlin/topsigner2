package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func NoCache() func(http.Handler) http.Handler {
	return middleware.NoCache
}
