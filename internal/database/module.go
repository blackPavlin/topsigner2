package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
)

var Module = fx.Module("database", fx.Provide(New), fx.Invoke(MakeMigrations))

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
}

type Result struct {
	fx.Out
	Pool *pgxpool.Pool
}

func New(params Params) (Result, error) {
	pool, err := NewPool(params.Config)
	if err != nil {
		return Result{}, fmt.Errorf("create database pool: %w", err)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := pool.Ping(ctx); err != nil {
				return fmt.Errorf("ping database: %w", err)
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			pool.Close()

			return nil
		},
	})

	return Result{Pool: pool}, nil
}
