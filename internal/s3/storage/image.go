package storage

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
	config *config.Config
	client *minio.Client
}

func NewImageStorage(config *config.Config, client *minio.Client) *ImageStorage {
	return &ImageStorage{
		config: config,
		client: client,
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
		return fmt.Errorf("upload image: %w", err)
	}

	return nil
}
