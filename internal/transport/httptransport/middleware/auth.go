package middleware

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/bboykiv/topsigner/gen/openapi"
)

func BearerAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if scopes := r.Context().Value(openapi.BearerAuthScopes); scopes == nil {
				next.ServeHTTP(w, r)
				
				return
			}

			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token == "" {
				render.Status(r, http.StatusUnauthorized)
				render.Respond(w, r, openapi.UnauthorizedJSONResponse{Message: "unauthorized"})

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
