package group

//go:generate go tool mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock -typed

import (
	"context"
	"time"

	"github.com/bboykiv/topsigner/internal/model"
)

type Repository interface {
	Get(ctx context.Context, filter *model.GroupFilter) (*model.Group, error)
	List(ctx context.Context, query *model.GroupQuery) ([]*model.Group, error)
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	Delete(ctx context.Context, filter *model.GroupFilter) error
}

type SessionRepository interface {
	Get(ctx context.Context, filter *model.SessionFilter) (*model.Session, error)
}

type VKClient interface {
	GetGroups(ctx context.Context, token string) ([]*Group, error)
	GenerateConnectGroupURL(groupID int64, state string) (string, error)
	ExchangeGroupCode(ctx context.Context, code string) (int64, string, error)
}

type StateRepository interface {
	Set(ctx context.Context, state string, userID int64, ttl time.Duration) error
	Pop(ctx context.Context, state string) (int64, error)
}
