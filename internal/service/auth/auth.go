package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

// todo: решить вопрос с сидированием и добавить создание первого пользователя

type UserRepository interface {
	Get(ctx context.Context, filter *model.UserFilter) (*model.User, error)
}

type SessionRepository interface{}

type Service struct {
	logger            *zap.Logger
	config            *config.Config
	userRepository    UserRepository
	sessionRepository SessionRepository
}

func New(logger *zap.Logger, config *config.Config, userRepository UserRepository) *Service {
	return &Service{
		logger:         logger,
		config:         config,
		userRepository: userRepository,
	}
}

func (s *Service) Login(ctx context.Context) (*model.AuthToken, error) {
	user, err := s.userRepository.Get(ctx, &model.UserFilter{})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, fmt.Errorf("get user: %w", err)
		}

		s.logger.Error("get user error", zap.Error(err))

		return nil, fmt.Errorf("get user: %w", err)
	}

	accessToken, err := s.SignToken(user)
	if err != nil {
		s.logger.Error("sign access token error", zap.Error(err))

		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		s.logger.Error("generate refresh token error", zap.Error(err))

		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &model.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *Service) Logout(ctx context.Context, userID int64) error {
	return nil
}

func (s *Service) Refresh(ctx context.Context) (*model.AuthToken, error) {
	return nil, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	token string,
) (*model.User, *model.Session, error) {
	claims, err := s.ParseAndValidateToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("parse and validate auth token: %w", err)
	}

	// todo: чтобы не грузить БД, имеет смысл добавить key/value хранилище
	user, err := s.userRepository.Get(ctx, &model.UserFilter{
		ID: model.IDFilter{Eq: new(claims.UserID)},
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, nil, model.ErrUserNotFound
		}

		s.logger.Error("get user error", zap.Error(err))

		return nil, nil, fmt.Errorf("get user: %w", err)
	}

	return user, &model.Session{}, nil
}

func (s *Service) SignToken(user *model.User) (string, error) {
	claims := &model.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	}

	// todo: возможно стоит вынести userID в subject

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(
		[]byte(s.config.Auth.SigningKey),
	)
	if err != nil {
		return "", fmt.Errorf("sign auth token string: %w", err)
	}

	return token, nil
}

func (s *Service) ParseAndValidateToken(token string) (*model.AuthClaims, error) {
	t, err := jwt.ParseWithClaims(token, &model.AuthClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}

		return []byte(s.config.Auth.SigningKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse auth claims: %w", err)
	}

	if claims, ok := t.Claims.(*model.AuthClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, model.ErrInvalidAuthToken
}
