package font

import (
	"context"

	"github.com/bboykiv/topsigner/internal/model"
)

type Repository interface {
	List(ctx context.Context, query *model.FontQuery) ([]*model.Font, error)
}
