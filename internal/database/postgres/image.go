package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bboykiv/topsigner/internal/model"
)

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{pool: pool}
}

func (r *ImageRepository) Create(ctx context.Context, image *model.Image) (*model.Image, error) {
	const query = `
		INSERT INTO images (name)
		VALUES ($1)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query, image.Name).Scan(
		&image.ID,
		&image.CreatedAt,
		&image.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, model.ErrImageAlreadyExists
			}
		}

		return nil, fmt.Errorf("create image: %w", err)
	}

	return image, nil
}

func (r *ImageRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
