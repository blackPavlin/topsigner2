package auth

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

type contextKey string

const userContextKey contextKey = "user"

func GetUserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)

	return user, ok
}

func SetUserToContext(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
