package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
)

var Module = fx.Module("s3", fx.Provide(New))

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
}

type Result struct {
	fx.Out
	Client *minio.Client
}

func New(params Params) (Result, error) {
	client, err := NewClient(params.Config)
	if err != nil {
		return Result{}, fmt.Errorf("create s3 client: %w", err)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := CreateBuckets(ctx, client, params.Config); err != nil {
				return fmt.Errorf("create buckets: %w", err)
			}

			return nil
		},
	})

	return Result{Client: client}, nil
}
