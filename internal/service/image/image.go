package image

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/model"
)

type Repository interface {
	List(ctx context.Context, query *model.ImageQuery) ([]*model.Image, error)
	Create(ctx context.Context, image *model.Image) (*model.Image, error)
	Delete(ctx context.Context, name string) error
}

type Storage interface {
	Upload(ctx context.Context, name string, reader io.Reader, size int64) error
	Delete(ctx context.Context, name string) error
}

type Service struct {
	logger     *zap.Logger
	repository Repository
	storage    Storage
}

func New(logger *zap.Logger, repository Repository, storage Storage) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
		storage:    storage,
	}
}

func (s *Service) List(
	ctx context.Context,
	query *model.ImageQuery,
) (*model.List[*model.Image], error) {
	pagination := model.Pagination{
		Cursor: query.Pagination.Cursor,
		Limit:  query.Pagination.Limit + 1,
	}

	images, err := s.repository.List(ctx, &model.ImageQuery{
		Filter:     query.Filter,
		Pagination: pagination,
	})
	if err != nil {
		s.logger.Error("get list of images error", zap.Error(err))

		return nil, fmt.Errorf("get list of images: %w", err)
	}

	result := &model.List[*model.Image]{
		Items: images,
	}

	if len(images) > query.Pagination.Limit {
		result.HasNext = true
		result.Items = images[:query.Pagination.Limit]

		lastItem := result.Items[len(result.Items)-1]
		result.NextCursor = new(model.EncodeCursor(lastItem.ID, lastItem.CreatedAt))

	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, fh *multipart.FileHeader) (*model.Image, error) {
	file, err := fh.Open()
	if err != nil {
		s.logger.Error("open multipart file error", zap.Error(err))

		return nil, fmt.Errorf("open multipart file: %w", err)
	}
	defer file.Close()

	if _, _, err := image.DecodeConfig(file); err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, model.ErrUnsupportedImageFormat
		}

		s.logger.Error("decode image config error", zap.Error(err))

		return nil, fmt.Errorf("decode image config: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		s.logger.Error("seek multipart file error", zap.Error(err))

		return nil, fmt.Errorf("seek multipart file: %w", err)
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fh.Filename))

	if err = s.storage.Upload(ctx, filename, file, fh.Size); err != nil {
		s.logger.Error("upload image to storage error", zap.Error(err))

		return nil, fmt.Errorf("upload image to storage: %w", err)
	}

	image := &model.Image{
		Name: filename,
	}

	if image, err = s.repository.Create(ctx, image); err != nil {
		s.logger.Error("create image error", zap.Error(err))

		return nil, fmt.Errorf("create image: %w", err)
	}

	return image, nil
}

func (s *Service) Delete(ctx context.Context, name string) error {
	if err := s.repository.Delete(ctx, name); err != nil {
		if errors.Is(err, model.ErrImageNotFound) {
			return nil
		}

		s.logger.Error("delete image error", zap.Error(err))

		return fmt.Errorf("delete image: %w", err)
	}

	if err := s.storage.Delete(ctx, name); err != nil {
		s.logger.Error("delete image from storage error", zap.Error(err))

		return fmt.Errorf("delete image from storage: %w", err)
	}

	return nil
}
