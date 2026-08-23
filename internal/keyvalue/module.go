package keyvalue

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
)

var Module = fx.Module("redis", fx.Provide(New))

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
}

type Result struct {
	fx.Out
	Client *redis.Client
}

func New(params Params) Result {
	client := NewClient(params.Config)

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.Ping(ctx).Err(); err != nil {
				return fmt.Errorf("ping redis: %w", err)
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := client.Close(); err != nil {
				return fmt.Errorf("close redis: %w", err)
			}

			return nil
		},
	})

	return Result{Client: client}
}
