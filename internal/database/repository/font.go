package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bboykiv/topsigner/internal/model"
)

type FontRepository struct {
	pool *pgxpool.Pool
}

func NewFontRepository(pool *pgxpool.Pool) *FontRepository {
	return &FontRepository{pool: pool}
}

func (r *FontRepository) List(
	ctx context.Context,
	query *model.FontQuery,
) ([]*model.Font, error) {
	builder := psql.Select("id", "name", "created_at", "updated_at").
		From(fontsTableName).
		OrderBy("created_at DESC", "id DESC").
		Limit(uint64(query.Pagination.Limit))

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

	fonts := make([]*model.Font, 0)

	for rows.Next() {
		font := &model.Font{}

		err := rows.Scan(&font.ID, &font.Name, &font.CreatedAt, &font.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan font: %w", err)
		}

		fonts = append(fonts, font)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate fonts: %w", err)
	}

	return fonts, nil
}
