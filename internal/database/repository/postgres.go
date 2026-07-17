package repository

import (
	"github.com/Masterminds/squirrel"

	"github.com/bboykiv/topsigner/internal/model"
)

const (
	userTableName    = "users"
	fontsTableName   = "fonts"
	imageTableName   = "images"
	sessionTableName = "sessions"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

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
