package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bboykiv/topsigner/internal/model"
)

type VKSessionRepository struct {
	pool *pgxpool.Pool
}

func NewVKSessionRepository(pool *pgxpool.Pool) *VKSessionRepository {
	return &VKSessionRepository{pool: pool}
}

func (r *VKSessionRepository) Get(
	ctx context.Context,
	filter *model.VKSessionFilter,
) (*model.VKSession, error) {
	builder := psql.Select(
		"id",
		"session_id",
		"device_id",
		"access_token_enc",
		"refresh_token_enc",
		"created_at",
		"updated_at",
	).
		From(vkSessionTable).
		Limit(1)

	builder = applyFilter(builder, "id", filter.ID)
	builder = applyFilter(builder, "session_id", filter.SessionID)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build sql query: %w", err)
	}

	vkSession := &model.VKSession{}

	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&vkSession.ID,
		&vkSession.SessionID,
		&vkSession.DeviceID,
		&vkSession.AccessTokenEnc,
		&vkSession.RefreshTokenEnc,
		&vkSession.CreatedAt,
		&vkSession.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrVKSessionNotFound
		}

		return nil, fmt.Errorf("get vk session: %w", err)
	}

	return vkSession, nil
}
