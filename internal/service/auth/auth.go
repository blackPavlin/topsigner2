package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

type UserRepository interface{}

type SessionRepository interface{}

type Service struct {
	logger            *zap.Logger
	config            *config.Config
	userRepository    UserRepository
	sessionRepository SessionRepository
}

func New(logger *zap.Logger, config *config.Config) *Service {
	return &Service{
		logger: logger,
		config: config,
	}
}

func (s *Service) Login(ctx context.Context) (*model.AuthToken, error) {
	accessToken, err := s.SignToken(777)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	return &model.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: "",
	}, nil
}

func (s *Service) Logout(ctx context.Context) error {
	return nil
}

func (s *Service) Refresh(ctx context.Context) error {
	return nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	token string,
) (*model.User, *model.Session, error) {
	claims, err := s.ParseAndValidateToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("parse and validate auth token: %w", err)
	}

	return &model.User{ID: claims.UserID}, &model.Session{}, nil
}

func (s *Service) SignToken(userID int64) (string, error) {
	claims := &model.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}

	// todo: возможно стоит вынести userID в subject
	// todo: возможно в методе стоит принимать user *model.User

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
