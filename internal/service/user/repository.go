package user

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

//go:generate go tool mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock -typed

type Repository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
}
