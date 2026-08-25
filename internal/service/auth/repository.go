package auth

//go:generate go tool mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock -typed

import (
	"context"
	"time"

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

type CodeVerifierRepository interface {
	Set(ctx context.Context, state, verifier string, ttl time.Duration) error
	Pop(ctx context.Context, state string) (string, error)
}

type UserCacheRepository interface {
	Get(ctx context.Context, userID int64) (*model.User, error)
	Set(ctx context.Context, user *model.User, ttl time.Duration) error
	Delete(ctx context.Context, userID int64) error
}

type VKIDClient interface {
	GenerateAuthorizationURL(challenge, state string) string
}
