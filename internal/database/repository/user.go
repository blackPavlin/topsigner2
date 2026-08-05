package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bboykiv/topsigner/internal/model"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Get(ctx context.Context, filter *model.UserFilter) (*model.User, error) {
	builder := psql.Select("id", "email", "password_hash", "role", "created_at", "updated_at").
		From(userTableName)

	builder = applyFilter(builder, "id", filter.ID)
	builder = applyFilter(builder, "email", filter.Email)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	user := &model.User{}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	sql, args, err := psql.Insert(userTableName).
		Columns("email", "password_hash", "role").
		Values(user.Email, user.PasswordHash, user.Role).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
