package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"

	"github.com/bboykiv/topsigner/internal/config"
)

func New(lc fx.Lifecycle, config *config.Config) (*minio.Client, error) {
	client, err := minio.New(config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.AccessKey, config.S3.SecretKey, ""),
		Region: config.S3.Region,
		Secure: config.S3.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			exists, err := client.BucketExists(ctx, config.S3.ImageBucket)
			if err != nil {
				return fmt.Errorf("check bucket %q exists: %w", config.S3.ImageBucket, err)
			}

			if !exists {
				err = client.MakeBucket(ctx, config.S3.ImageBucket, minio.MakeBucketOptions{
					Region: config.S3.Region,
				})
				if err != nil {
					return fmt.Errorf("create bucket %q: %w", config.S3.ImageBucket, err)
				}
			}

			// todo: иницализация font_bucket

			return nil
		},
	})

	return client, nil
}
