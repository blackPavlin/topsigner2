package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/font"
	"github.com/bboykiv/topsigner/internal/service/group"
	"github.com/bboykiv/topsigner/internal/service/health"
	"github.com/bboykiv/topsigner/internal/service/image"
	"github.com/bboykiv/topsigner/internal/service/user"
)

var Module = fx.Module("postgres",
	fx.Provide(New),
	fx.Provide(
		fx.Annotate(
			NewGroupRepository,
			fx.As(new(group.Repository)),
		),
		fx.Annotate(
			NewFontRepository,
			fx.As(new(font.Repository)),
		),
		fx.Annotate(
			NewImageRepository,
			fx.As(new(image.Repository)),
		),
		fx.Annotate(
			NewUserRepository,
			fx.As(new(user.Repository)),
			fx.As(new(auth.UserRepository)),
		),
		fx.Annotate(
			NewSessionRepository,
			fx.As(new(auth.SessionRepository)),
			fx.As(new(group.SessionRepository)),
		),
		fx.Annotate(
			NewHealth,
			fx.As(new(health.DatabaseChecker)),
		),
	),
	fx.Invoke(MakeMigrations),
)

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
