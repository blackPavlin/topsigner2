package user

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger     *zap.Logger
	config     *config.Config
	repository Repository
}

func New(logger *zap.Logger, config *config.Config, repository Repository) *Service {
	return &Service{
		logger:     logger.Named("user-service"),
		config:     config,
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, input *CreateUserInput) (*model.User, error) {
	passwordHash, err := model.GeneratePasswordHash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("generate user password hash: %w", err)
	}

	user := &model.User{
		Email:        new(input.Email),
		PasswordHash: new(passwordHash),
		Role:         input.Role,
	}

	if user, err = s.repository.Create(ctx, user); err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, model.ErrUserAlreadyExists
		}

		s.logger.Error("create user", zap.Error(err))

		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *Service) CreateDefault(ctx context.Context) error {
	input := &CreateUserInput{
		Email:    s.config.User.Default.Email,
		Password: s.config.User.Default.Password,
		Role:     s.config.User.Default.Role,
	}

	if _, err := s.Create(ctx, input); err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			s.logger.Info("default user already exists")

			return nil
		}

		s.logger.Error("create default user", zap.Error(err))

		return fmt.Errorf("create default user: %w", err)
	}

	s.logger.Info("default user successful create")

	return nil
}
