package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func IP() func(http.Handler) http.Handler {
	// todo: ClientIPFromRemoteAddr не будет корректно работать если сервер стоит за reverse proxy
	return middleware.ClientIPFromRemoteAddr
}
