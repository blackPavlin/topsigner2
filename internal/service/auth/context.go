package auth

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

type contextKey string

const (
	userContextKey    contextKey = "user"
	sessionContextKey contextKey = "session"
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
