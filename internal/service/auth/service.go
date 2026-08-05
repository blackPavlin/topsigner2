package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

// todo: решить вопрос с сидированием и добавить создание первого пользователя
// todo: чтобы не грузить БД, имеет смысл добавить key/value хранилище
// todo: возвращать в TokenPair ExpiresAt

type Service struct {
	logger            *zap.Logger
	config            *config.Config
	vkidClient        VKIDClient
	userRepository    UserRepository
	sessionRepository SessionRepository
}

func New(
	logger *zap.Logger,
	config *config.Config,
	vkidClient VKIDClient,
	userRepository UserRepository,
	sessionRepository SessionRepository,
) *Service {
	return &Service{
		logger:            logger,
		config:            config,
		vkidClient:        vkidClient,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

func (s *Service) Login(ctx context.Context, input *LoginInput) (*TokenPair, error) {
	user, err := s.userRepository.Get(ctx, &model.UserFilter{
		Email: model.TextFilter{Eq: new(input.Email)},
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, fmt.Errorf("get user: %w", err)
		}

		s.logger.Error("get user error", zap.Error(err))

		return nil, fmt.Errorf("get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		s.logger.Error("generate refresh token error", zap.Error(err))

		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session := &model.Session{
		UserID:           user.ID,
		IP:               input.IP,
		UserAgent:        input.UserAgent,
		RefreshTokenHash: hashRefreshToken(refreshToken),
		ExpiresAt:        time.Now().Add(s.config.Auth.RefreshTokenTTL),
	}

	// todo: не создавать новую сессию на каждый запрос авторизации, если user_id, ip и user_agent совпадают

	session, err = s.sessionRepository.Create(ctx, session)
	if err != nil {
		s.logger.Error("create session error", zap.Error(err))

		return nil, fmt.Errorf("create session: %w", err)
	}

	accessToken, err := s.SignAccessToken(user.ID, session.ID)
	if err != nil {
		s.logger.Error("sign access token error", zap.Error(err))

		return nil, fmt.Errorf("sign access token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, userID int64, refreshToken *string) error {
	filter := &model.SessionFilter{
		UserID: model.IDFilter{Eq: new(userID)},
	}

	if refreshToken != nil {
		filter.RefreshTokenHash = model.TextFilter{
			Eq: new(hashRefreshToken(*refreshToken)),
		}
	}

	if err := s.sessionRepository.Delete(ctx, filter); err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil
		}

		s.logger.Error("delete session error", zap.Error(err))

		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	session, err := s.sessionRepository.Get(ctx, &model.SessionFilter{
		RefreshTokenHash: model.TextFilter{
			Eq: new(hashRefreshToken(refreshToken)),
		},
	})
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, fmt.Errorf("get session: %w", err)
		}

		s.logger.Error("get session error", zap.Error(err))

		return nil, fmt.Errorf("get session: %w", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		err = s.sessionRepository.Delete(ctx, &model.SessionFilter{
			UserID:           model.IDFilter{Eq: new(session.UserID)},
			RefreshTokenHash: model.TextFilter{Eq: new(session.RefreshTokenHash)},
		})
		if err != nil {
			if errors.Is(err, model.ErrSessionNotFound) {
				return nil, fmt.Errorf("delete session: %w", err)
			}

			s.logger.Error("delete session error", zap.Error(err))

			return nil, fmt.Errorf("delete session: %w", err)
		}

		return nil, ErrTokenIsExpired
	}

	refreshToken, err = generateRefreshToken()
	if err != nil {
		s.logger.Error("generate refresh token error", zap.Error(err))

		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session.RefreshTokenHash = hashRefreshToken(refreshToken)
	session.ExpiresAt = time.Now().Add(s.config.Auth.RefreshTokenTTL)

	// todo: решить нужно ли обновлять user_agent и ip

	session, err = s.sessionRepository.Update(ctx, session)
	if err != nil {
		s.logger.Error("update session error", zap.Error(err))

		return nil, fmt.Errorf("update session: %w", err)
	}

	accessToken, err := s.SignAccessToken(session.UserID, session.ID)
	if err != nil {
		s.logger.Error("sign access token error", zap.Error(err))

		return nil, fmt.Errorf("sign access token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Authorize(ctx context.Context, token string) (*model.User, error) {
	claims, err := s.ParseAndValidateAccessToken(token)
	if err != nil {
		return nil, fmt.Errorf("parse and validate auth token: %w", err)
	}

	user, err := s.userRepository.Get(ctx, &model.UserFilter{
		ID: model.IDFilter{Eq: new(claims.UserID)},
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}

		s.logger.Error("get user error", zap.Error(err))

		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *Service) SignAccessToken(userID int64, sessionID string) (string, error) {
	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:    userID,
		SessionID: sessionID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(
		[]byte(s.config.Auth.SigningKey),
	)
	if err != nil {
		return "", fmt.Errorf("sign auth token string: %w", err)
	}

	return token, nil
}

func (s *Service) ParseAndValidateAccessToken(token string) (*AccessTokenClaims, error) {
	t, err := jwt.ParseWithClaims(token, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}

		return []byte(s.config.Auth.SigningKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse auth claims: %w", err)
	}

	if claims, ok := t.Claims.(*AccessTokenClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, ErrInvalidAuthToken
}

func (s *Service) GenerateVKIDAuthorizationURL() (string, error) {
	codeVerifier, err := generateRandomString(codeVerifierBytes)
	if err != nil {
		s.logger.Error("generate code verifier error", zap.Error(err))

		return "", fmt.Errorf("generate code verifier: %w", err)
	}

	state, err := generateRandomString(stateBytes)
	if err != nil {
		s.logger.Error("generate state error", zap.Error(err))

		return "", fmt.Errorf("generate state: %w", err)
	}

	// todo: после добавления key/value хранилища, записывать codeVerifier в него

	return s.vkidClient.GenerateAuthorizationURL(codeChallengeS256(codeVerifier), state), nil
}
