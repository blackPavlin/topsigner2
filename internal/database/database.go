package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
)

func New(lc fx.Lifecycle, config *config.Config) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(config.Postgres.ToDataSource())
	if err != nil {
		return nil, fmt.Errorf("parse database pool config: %w", err)
	}

	conf.MaxConns = config.Postgres.MaxOpenConns
	conf.MinConns = config.Postgres.MaxIdleConns
	conf.MaxConnLifetime = config.Postgres.ConnMaxLifetime
	conf.MaxConnIdleTime = config.Postgres.ConnMaxIdleTime
	conf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err = pool.Ping(ctx); err != nil {
				return fmt.Errorf("ping database: %w", err)
			}

			return nil
		},
		OnStop: func(_ context.Context) error {
			pool.Close()

			return nil
		},
	})

	return pool, nil
}
