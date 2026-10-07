package middleware

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/bboykiv/topsigner/internal/model"
)

type contextKey string

const (
	userContextKey      contextKey = "user"
	sessionIDContextKey contextKey = "session_id"
	userAgentContextKey contextKey = "user-agent"
)

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)

	return user, ok
}

func SetUserToContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetSessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(sessionIDContextKey).(string)

	return sessionID, ok
}

func SetSessionIDToContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDContextKey, sessionID)
}

func GetUserAgentFromContext(ctx context.Context) string {
	userAgent, _ := ctx.Value(userAgentContextKey).(string)

	return userAgent
}

func SetUserAgentToContext(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, userAgentContextKey, userAgent)
}

func GetClientIPFromContext(ctx context.Context) string {
	return middleware.GetClientIP(ctx)
}
