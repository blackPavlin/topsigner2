package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Health struct {
	client *redis.Client
}

func NewHealth(client *redis.Client) *Health {
	return &Health{client: client}
}

func (h *Health) Ping(ctx context.Context) error {
	if err := h.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis cache: %w", err)
	}

	return nil
}
