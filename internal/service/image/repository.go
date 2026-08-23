package image

//go:generate go tool mockgen -source=repository.go -destination=mock/mock_repository.go -package=mock -typed

import (
	"context"
	"io"

	"github.com/bboykiv/topsigner/internal/model"
)

type Repository interface {
	List(ctx context.Context, query *model.ImageQuery) ([]*model.Image, error)
	Create(ctx context.Context, image *model.Image) (*model.Image, error)
	Delete(ctx context.Context, userID int64, name string) error
}

type Storage interface {
	Upload(ctx context.Context, name string, reader io.Reader, size int64) error
	Delete(ctx context.Context, name string) error
}
