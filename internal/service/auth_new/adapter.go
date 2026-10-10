package auth

//go:generate go tool mockgen -source=adapter.go -destination=mock/mock_adapter.go -package=mock -typed

import (
	"context"
	"time"
	"uuid"

	"github.com/bboykiv/topsigner/internal/model"
)

type UserRepository interface {
	Get(ctx context.Context, filter *model.UserFilter) (*model.User, error)
}

type UserCache interface {
	Get(ctx context.Context, userID int64) (*model.User, error)
}

type PasswordHasher interface {
	Compare(hash, password string) error
}

type TokenIssuer interface {
	Issue(ttl time.Duration) (string, error)
}

type TokenGenerator interface {
	New() (string, error)
	Hash(token string) string
}

type SessionRepository interface {
	Get(ctx context.Context, filter *model.SessionFilter) (*model.Session, error)
	Create(ctx context.Context, session *model.Session) (*model.Session, error)
}

type SessionCache interface {
	Get(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
}
