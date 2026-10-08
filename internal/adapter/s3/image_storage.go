package s3

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	"github.com/bboykiv/topsigner/internal/config"
)

type ImageStorage struct {
	client *minio.Client
	config *config.Config
}

func NewImageStorage(client *minio.Client, config *config.Config) *ImageStorage {
	return &ImageStorage{
		client: client,
		config: config,
	}
}

func (s *ImageStorage) Upload(ctx context.Context, name string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx,
		s.config.S3.ImageBucket,
		name,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: mime.TypeByExtension(filepath.Ext(name)),
		},
	)
	if err != nil {
		return fmt.Errorf("upload image object: %w", err)
	}

	return nil
}

func (s *ImageStorage) Delete(ctx context.Context, name string) error {
	err := s.client.RemoveObject(ctx, s.config.S3.ImageBucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("delete image object: %w", err)
	}

	return nil
}
