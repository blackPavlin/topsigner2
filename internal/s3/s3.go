package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/bboykiv/topsigner/internal/config"
)

func NewClient(config *config.Config) (*minio.Client, error) {
	client, err := minio.New(config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.AccessKey, config.S3.SecretKey, ""),
		Region: config.S3.Region,
		Secure: config.S3.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return client, nil
}

func CreateBukets(ctx context.Context, client *minio.Client, config *config.Config) error {
	buckets := []string{config.S3.ImageBucket, config.S3.FontBucket}

	for _, bucket := range buckets {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("check bucket %q exists: %w", bucket, err)
		}

		if !exists {
			err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
				Region: config.S3.Region,
			})
			if err != nil {
				return fmt.Errorf("create bucket %q: %w", bucket, err)
			}
		}
	}

	return nil
}
