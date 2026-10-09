package image

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"uuid"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger     *zap.Logger
	repository Repository
	storage    Storage
}

func New(logger *zap.Logger, repository Repository, storage Storage) *Service {
	return &Service{
		logger:     logger.Named("image-service"),
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

func (s *Service) Create(ctx context.Context, userID int64, upload *Upload) (*model.Image, error) {
	_, format, err := image.DecodeConfig(upload.File)
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, model.ErrUnsupportedImageFormat
		}

		s.logger.Error("decode image config error", zap.Error(err))

		return nil, fmt.Errorf("decode image config: %w", err)
	}

	if _, err := upload.File.Seek(0, io.SeekStart); err != nil {
		s.logger.Error("seek image file error", zap.Error(err))

		return nil, fmt.Errorf("seek image file: %w", err)
	}

	filename := fmt.Sprintf("%s.%s", uuid.New().String(), format)

	if err := s.storage.Upload(ctx, filename, upload.File, upload.Size); err != nil {
		s.logger.Error("upload image to storage error", zap.Error(err))

		return nil, fmt.Errorf("upload image to storage: %w", err)
	}

	image := &model.Image{
		Name:   filename,
		UserID: userID,
	}

	if image, err = s.repository.Create(ctx, image); err != nil {
		s.logger.Error("create image error", zap.Error(err))

		return nil, fmt.Errorf("create image: %w", err)
	}

	return image, nil
}

func (s *Service) Delete(ctx context.Context, imageID, userID int64) error {
	image, err := s.repository.Delete(ctx, imageID, userID)
	if  err != nil {
		if errors.Is(err, model.ErrImageNotFound) {
			return model.ErrImageNotFound
		}

		s.logger.Error("delete user image", zap.Error(err))

		return fmt.Errorf("delete user image: %w", err)
	}

	if err = s.storage.Delete(ctx, image.Name); err != nil {
		s.logger.Error("delete user image from storage", zap.Error(err))

		return fmt.Errorf("delete user image from storage: %w", err)
	}

	return nil
}
