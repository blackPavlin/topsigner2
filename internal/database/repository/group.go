package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bboykiv/topsigner/internal/model"
)

type GroupRepository struct {
	pool *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{pool: pool}
}

func (r *GroupRepository) Get(
	ctx context.Context,
	filter *model.GroupFilter,
) (*model.Group, error) {
	builder := psql.Select("id", "user_id", "external_id", "created_at", "updated_at").
		From(groupTableName).
		Limit(1)

	builder = applyFilter(builder, "id", filter.ID)
	builder = applyFilter(builder, "user_id", filter.UserID)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	group := &model.Group{}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&group.ID,
		&group.UserID,
		&group.ExternalID,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrGroupNotFound
		}

		return nil, fmt.Errorf("get group: %w", err)
	}

	return group, nil
}

func (r *GroupRepository) List(
	ctx context.Context,
	query *model.GroupQuery,
) ([]*model.Group, error) {
	builder := psql.Select("id", "user_id", "external_id", "created_at", "updated_at").
		From(groupTableName).
		OrderBy("created_at DESC", "id DESC").
		Limit(uint64(query.Pagination.Limit))

	builder = applyFilter(builder, "user_id", query.Filter.UserID)

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
		return nil, fmt.Errorf("query groups: %w", err)
	}
	defer rows.Close()

	groups := make([]*model.Group, 0)

	for rows.Next() {
		group := &model.Group{}

		err := rows.Scan(
			&group.ID,
			&group.UserID,
			&group.ExternalID,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate groups: %w", err)
	}

	return groups, nil
}

func (r *GroupRepository) Create(ctx context.Context, group *model.Group) (*model.Group, error) {
	sql, args, err := psql.Insert(groupTableName).
		Columns("user_id", "external_id").
		Values(group.UserID).
		Suffix("ON CONFLICT DO NOTHING RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&group.ID,
		&group.ExternalID,
		&group.CreatedAt,
		&group.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrGroupAlreadyExists
		}

		return nil, fmt.Errorf("create group: %w", err)
	}

	return nil, nil
}

func (r *GroupRepository) Delete(ctx context.Context, filter *model.GroupFilter) error {
	if filter.UserID.Eq == nil || filter.ID.Eq == nil {
		return nil
	}

	sql, args, err := psql.Delete(groupTableName).
		Where(squirrel.Eq{"user_id": filter.UserID.Eq, "id": filter.ID.Eq}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build sql query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete group: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrGroupNotFound
	}

	return nil
}
