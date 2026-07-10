package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/bboykiv/topsigner/gen/openapi"
	"github.com/bboykiv/topsigner/internal/model"
)

type AuthService interface {
	Authenticate(ctx context.Context, token string) (*model.User, *model.Session, error)
}

func BearerAuth(authService AuthService) func(next http.Handler) http.Handler {
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

			ctx := r.Context()

			user, session, err := authService.Authenticate(ctx, token)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.Respond(w, r, openapi.UnauthorizedJSONResponse{Message: "unauthorized"})

				return
			}

			ctx = model.SetUserToContext(ctx, user)
			ctx = model.SetSessionToContext(ctx, session)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
