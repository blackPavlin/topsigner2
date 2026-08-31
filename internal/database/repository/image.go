package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
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
	builder := psql.Select("id", "user_id", "name", "created_at", "updated_at").
		From(imageTableName).
		OrderBy("created_at DESC", "id DESC").
		Limit(uint64(query.Pagination.Limit))

	builder = applyFilter(builder, "user_id", query.Filter.UserID)
	builder = applyFilter(builder, "name", query.Filter.Name)

	if cursor := query.Pagination.Cursor; cursor != nil {
		builder = builder.Where(
			squirrel.Expr("(id, created_at) < (?, ?)", cursor.ID, cursor.CreatedAt),
		)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query images: %w", err)
	}
	defer rows.Close()

	images := make([]*model.Image, 0)

	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(&image.ID, &image.UserID, &image.Name, &image.CreatedAt, &image.UpdatedAt)
		if err != nil {
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
	sql, args, err := psql.Insert(imageTableName).
		Columns("user_id", "name").
		Values(image.UserID, image.Name).
		Suffix("ON CONFLICT DO NOTHING RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&image.ID,
		&image.CreatedAt,
		&image.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrImageAlreadyExists
		}

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return nil, model.ErrUserNotFound
			}
		}

		return nil, fmt.Errorf("create image: %w", err)
	}

	return image, nil
}

func (r *ImageRepository) Delete(ctx context.Context, userID int64, name string) error {
	sql, args, err := psql.Delete(imageTableName).
		Where(squirrel.Eq{"user_id": userID, "name": name}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrImageNotFound
	}

	return nil
}
