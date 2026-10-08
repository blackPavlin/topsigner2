package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Health struct {
	pool *pgxpool.Pool
}

func NewHealth(pool *pgxpool.Pool) *Health {
	return &Health{pool: pool}
}

func (h *Health) Ping(ctx context.Context) error {
	if err := h.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres database: %w", err)
	}

	return nil
}
