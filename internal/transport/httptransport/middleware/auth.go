package middleware

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/bboykiv/topsigner/gen/httpserver"
	"github.com/bboykiv/topsigner/internal/service/auth"
)

func BearerAuth(authService *auth.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if scopes := r.Context().Value(httpserver.BearerAuthScopes); scopes == nil {
				next.ServeHTTP(w, r)

				return
			}

			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if token == "" {
				render.Status(r, http.StatusUnauthorized)
				render.Respond(w, r, httpserver.UnauthorizedJSONResponse{Message: "unauthorized"})

				return
			}

			ctx := r.Context()

			user, err := authService.Authorize(ctx, token)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.Respond(w, r, httpserver.UnauthorizedJSONResponse{Message: "unauthorized"})

				return
			}

			ctx = auth.SetUserToContext(ctx, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
