package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

const (
	userTableName    = "users"
	fontsTableName   = "fonts"
	imageTableName   = "images"
	groupTableName   = "groups"
	sessionTableName = "sessions"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewPool(config *config.Config) (*pgxpool.Pool, error) {
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

	return pool, nil
}

func applyFilter[T comparable](
	builder squirrel.SelectBuilder,
	column string,
	filter model.Filter[T],
) squirrel.SelectBuilder {
	if filter.Eq != nil {
		builder = builder.Where(squirrel.Eq{column: *filter.Eq})
	}

	if filter.Neq != nil {
		builder = builder.Where(squirrel.NotEq{column: *filter.Neq})
	}

	if len(filter.In) > 0 {
		builder = builder.Where(squirrel.Eq{column: filter.In})
	}

	if len(filter.NotIn) > 0 {
		builder = builder.Where(squirrel.NotEq{column: filter.NotIn})
	}

	return builder
}
