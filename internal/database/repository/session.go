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
			"auth_type",
			"ip",
			"user_agent",
			"refresh_token_hash",
			"oauth_device_id",
			"oauth_access_token_enc",
			"oauth_refresh_token_enc",
			"expires_at",
			"created_at",
			"updated_at",
		).
		From(sessionTableName).
		Limit(1)

	builder = applyFilter(builder, "id", filter.ID)
	builder = applyFilter(builder, "user_id", filter.UserID)
	builder = applyFilter(builder, "auth_type", filter.AuthType)
	builder = applyFilter(builder, "refresh_token_hash", filter.RefreshTokenHash)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	session := &model.Session{}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&session.ID,
		&session.UserID,
		&session.AuthType,
		&session.IP,
		&session.UserAgent,
		&session.RefreshTokenHash,
		&session.OAuthDeviceID,
		&session.OAuthAccessTokenEnc,
		&session.OAuthRefreshTokenEnc,
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

func (r *SessionRepository) List(
	ctx context.Context,
	query *model.SessionQuery,
) ([]*model.Session, error) {
	builder := psql.Select(
		"id::text",
		"user_id",
		"auth_type",
		"ip",
		"user_agent",
		"refresh_token_hash",
		"oauth_device_id",
		"oauth_access_token_enc",
		"oauth_refresh_token_enc",
		"expires_at",
		"created_at",
		"updated_at",
	).
		From(sessionTableName)

	builder = applyFilter(builder, "id", query.Filter.ID)
	builder = applyFilter(builder, "user_id", query.Filter.UserID)
	builder = applyFilter(builder, "auth_type", query.Filter.AuthType)
	builder = applyFilter(builder, "refresh_token_hash", query.Filter.RefreshTokenHash)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]*model.Session, 0)

	for rows.Next() {
		session := &model.Session{}

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.AuthType,
			&session.IP,
			&session.UserAgent,
			&session.RefreshTokenHash,
			&session.OAuthDeviceID,
			&session.OAuthAccessTokenEnc,
			&session.OAuthRefreshTokenEnc,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepository) Create(
	ctx context.Context,
	session *model.Session,
) (*model.Session, error) {
	sql, args, err := psql.Insert(sessionTableName).
		Columns(
			"user_id",
			"auth_type",
			"ip",
			"user_agent",
			"refresh_token_hash",
			"oauth_device_id",
			"oauth_access_token_enc",
			"oauth_refresh_token_enc",
			"expires_at",
		).
		Values(
			session.UserID,
			session.AuthType,
			session.IP,
			session.UserAgent,
			session.RefreshTokenHash,
			session.OAuthDeviceID,
			session.OAuthAccessTokenEnc,
			session.OAuthRefreshTokenEnc,
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
		Set("oauth_access_token_enc", session.OAuthAccessTokenEnc).
		Set("oauth_refresh_token_enc", session.OAuthRefreshTokenEnc).
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
		return nil
	}

	builder := psql.Delete(sessionTableName).
		Where(squirrel.Eq{"user_id": *filter.UserID.Eq})

	if filter.ID.Eq != nil {
		builder = builder.Where(squirrel.Eq{"id": *filter.ID.Eq})
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
