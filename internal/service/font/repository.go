package font

//go:generate go tool mockgen -source=repository.go -destination=mocks/mock_repository.go -package=mocks -typed

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

type Repository interface {
	List(ctx context.Context, query *model.FontQuery) ([]*model.Font, error)
}
