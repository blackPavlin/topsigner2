package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *ImageRepository) List(
	ctx context.Context,
	query *model.ImageQuery,
) ([]*model.Image, error) {
	const q = `
		SELECT id, name, created_at, updated_at
		FROM images
		WHERE ($1::bigint IS NULL AND $2::timestamptz IS NULL)
		   OR (id, created_at) < ($1, $2) 
		ORDER BY created_at DESC, id ASC
		LIMIT $3`

	var (
		cursorID        *int64
		cursorCreatedAt *time.Time
	)

	if query.Pagination.Cursor != nil {
		cursorID = &query.Pagination.Cursor.ID
		cursorCreatedAt = &query.Pagination.Cursor.CreatedAt
	}

	rows, err := r.pool.Query(ctx, q, cursorID, cursorCreatedAt, query.Pagination.Limit)
	if err != nil {
		return nil, fmt.Errorf("query images: %w", err)
	}
	defer rows.Close()

	images := make([]*model.Image, 0)

	for rows.Next() {
		image := &model.Image{}

		if err := rows.Scan(&image.ID, &image.Name, &image.CreatedAt, &image.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan image: %w", err)
		}

		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate images: %w", err)
	}

	return images, nil
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

func (r *ImageRepository) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM images WHERE name = $1`

	tag, err := r.pool.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrImageNotFound
	}

	return nil
}
