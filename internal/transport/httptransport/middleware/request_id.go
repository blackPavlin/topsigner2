package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestID() func(http.Handler) http.Handler {
	return middleware.RequestID
}
