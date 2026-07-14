package auth

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

type UserRepository interface {
	Get(ctx context.Context, filter *model.UserFilter) (*model.User, error)
}

type SessionRepository interface {
	Get(ctx context.Context, filter *model.SessionFilter) (*model.Session, error)
	Create(ctx context.Context, session *model.Session) (*model.Session, error)
	Update(ctx context.Context, session *model.Session) (*model.Session, error)
	Delete(ctx context.Context, filter *model.SessionFilter) error
}
