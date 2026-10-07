package middleware

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/bboykiv/topsigner/internal/model"
)

type contextKey string

const (
	userContextKey      contextKey = "user"
	sessionContextKey   contextKey = "session"
	userAgentContextKey contextKey = "user-agent"
)

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)

	return user, ok
}

func SetUserToContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetSessionFromContext(ctx context.Context) (*model.Session, bool) {
	session, ok := ctx.Value(sessionContextKey).(*model.Session)

	return session, ok
}

func SetSessionToContext(ctx context.Context, session *model.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
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
