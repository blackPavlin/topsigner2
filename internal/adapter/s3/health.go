package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	"github.com/bboykiv/topsigner/internal/config"
)

type Health struct {
	client *minio.Client
	config *config.Config
}

func NewHealth(client *minio.Client, config *config.Config) *Health {
	return &Health{client: client, config: config}
}

func (h *Health) Ping(ctx context.Context) error {
	buckets := []string{h.config.S3.ImageBucket, h.config.S3.FontBucket}

	// todo: распараллелить проверку бакетов

	for _, bucket := range buckets {
		exists, err := h.client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("check bucket %q exists: %w", bucket, err)
		}

		if !exists {
			return fmt.Errorf("bucket %q not exists", bucket)
		}
	}

	return nil
}
