package font

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger     *zap.Logger
	repository Repository
}

func New(logger *zap.Logger, repository Repository) *Service {
	return &Service{
		logger:     logger.Named("font-service"),
		repository: repository,
	}
}

func (s *Service) List(
	ctx context.Context,
	query *model.FontQuery,
) (*model.List[*model.Font], error) {
	pagination := model.Pagination{
		Cursor: query.Pagination.Cursor,
		Limit:  query.Pagination.Limit + 1,
	}

	fonts, err := s.repository.List(ctx, &model.FontQuery{
		Filter:     query.Filter,
		Pagination: pagination,
	})
	if err != nil {
		s.logger.Error("get list of fonts error", zap.Error(err))

		return nil, fmt.Errorf("get list of fonts: %w", err)
	}

	result := &model.List[*model.Font]{
		Items: fonts,
	}

	if len(fonts) > query.Pagination.Limit {
		result.HasNext = true
		result.Items = fonts[:query.Pagination.Limit]

		lastItem := result.Items[len(result.Items)-1]
		result.NextCursor = new(model.EncodeCursor(lastItem.ID, lastItem.CreatedAt))
	}

	return result, nil
}
