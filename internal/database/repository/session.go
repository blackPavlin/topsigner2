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

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Get(
	ctx context.Context,
	filter *model.SessionFilter,
) (*model.Session, error) {
	builder := psql.
		Select(
			"id::text",
			"user_id",
			"refresh_token_hash",
			"expires_at",
			"created_at",
			"updated_at",
		).
		From(sessionTableName)

	builder = applyFilter(builder, "id", filter.ID)
	builder = applyFilter(builder, "user_id", filter.UserID)
	builder = applyFilter(builder, "refresh_token_hash", filter.RefreshTokenHash)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	session := &model.Session{}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrSessionNotFound
		}

		return nil, fmt.Errorf("get session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) Create(
	ctx context.Context,
	session *model.Session,
) (*model.Session, error) {
	sql, args, err := psql.Insert(sessionTableName).
		Columns(
			"user_id",
			"refresh_token_hash",
			"expires_at",
		).
		Values(
			session.UserID,
			session.RefreshTokenHash,
			session.ExpiresAt,
		).
		Suffix("RETURNING id::text, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return nil, model.ErrUserNotFound
			}
		}

		return nil, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) Update(
	ctx context.Context,
	session *model.Session,
) (*model.Session, error) {
	sql, args, err := psql.Update(sessionTableName).
		Set("refresh_token_hash", session.RefreshTokenHash).
		Set("expires_at", session.ExpiresAt).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": session.ID}).
		Suffix("RETURNING id::text, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&session.ID,
		&session.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, filter *model.SessionFilter) error {
	if filter.UserID.Eq == nil {
		return model.ErrSessionNotFound
	}

	builder := psql.Delete(sessionTableName).
		Where(squirrel.Eq{"user_id": *filter.UserID.Eq})

	if filter.RefreshTokenHash.Eq != nil {
		builder = builder.Where(squirrel.Eq{"refresh_token_hash": *filter.RefreshTokenHash.Eq})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build sql query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSessionNotFound
	}

	return nil
}
