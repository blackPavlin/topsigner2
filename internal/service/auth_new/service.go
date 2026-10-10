package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger            *zap.Logger
	userRepository    UserRepository
	sessionRepository SessionRepository
	passwordHasher    PasswordHasher
}

func NewService(
	logger *zap.Logger,
	userRepository UserRepository,
	sessionRepository SessionRepository,
	passwordHasher PasswordHasher,
) *Service {
	return &Service{
		logger:            logger.Named("auth-service"),
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		passwordHasher:    passwordHasher,
	}
}

func (s *Service) Login(ctx context.Context, login *Login) (*Token, error) {
	user, err := s.userRepository.Get(ctx, &model.UserFilter{
		Email: model.TextFilter{Eq: new(strings.ToLower(login.Email))},
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}

		s.logger.Error("get user", zap.Error(err))

		return nil, fmt.Errorf("get user: %w", err)
	}

	if user.PasswordHash == nil {
		return nil, model.ErrPasswordLoginNotAvailable
	}

	if err = s.passwordHasher.Compare(*user.PasswordHash, login.Password); err != nil {
		return nil, model.ErrInvalidPassword
	}

	// todo: сгенерировать refresh_token
	// todo: создать сессию
	// todo: подписать jwt токен
	// todo: положить в кэш user и session
}
