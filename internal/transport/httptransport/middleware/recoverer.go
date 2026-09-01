package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Recoverer() func(http.Handler) http.Handler {
	return middleware.Recoverer
}
