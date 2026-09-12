package group

//go:generate go tool mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock -typed

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

type Repository interface {
	Get(ctx context.Context, filter *model.GroupFilter) (*model.Group, error)
	List(ctx context.Context, query *model.GroupQuery) ([]*model.Group, error)
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	Delete(ctx context.Context, filter *model.GroupFilter) error
}

type VKClient interface{}
