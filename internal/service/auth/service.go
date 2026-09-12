package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger                 *zap.Logger
	config                 *config.Config
	encryptor              *Encryptor
	vkidClient             VKIDClient
	userRepository         UserRepository
	sessionRepository      SessionRepository
	userCacheRepository    UserCacheRepository
	sessionCacheRepository SessionCacheRepository
	codeVerifierRepository CodeVerifierRepository
}

func New(
	logger *zap.Logger,
	config *config.Config,
	vkidClient VKIDClient,
	userRepository UserRepository,
	sessionRepository SessionRepository,
	userCacheRepository UserCacheRepository,
	sessionCacheRepository SessionCacheRepository,
	codeVerifierRepository CodeVerifierRepository,
) (*Service, error) {
	encryptor, err := NewEncryptor(config.Auth.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("create new encryptor: %w", err)
	}

	return &Service{
		logger:                 logger.Named("auth-service"),
		config:                 config,
		encryptor:              encryptor,
		vkidClient:             vkidClient,
		userRepository:         userRepository,
		sessionRepository:      sessionRepository,
		userCacheRepository:    userCacheRepository,
		sessionCacheRepository: sessionCacheRepository,
		codeVerifierRepository: codeVerifierRepository,
	}, nil
}

func (s *Service) Login(ctx context.Context, input *LoginInput) (*TokenPair, error) {
	user, err := s.userRepository.Get(ctx, &model.UserFilter{
		Email: model.TextFilter{Eq: new(input.Email)},
	})
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}

		s.logger.Error("get user", zap.Error(err))

		return nil, fmt.Errorf("get user: %w", err)
	}

	if user.PasswordHash == nil {
		return nil, ErrPasswordLoginNotAvailable
	}

	if err = model.ComparePasswordAndHash(*user.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		s.logger.Error("generate refresh token", zap.Error(err))

		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session := &model.Session{
		UserID:           user.ID,
		AuthType:         model.AuthTypePassword,
		IP:               input.IP,
		UserAgent:        input.UserAgent,
		RefreshTokenHash: hashRefreshToken(refreshToken),
		ExpiresAt:        time.Now().Add(s.config.Auth.RefreshTokenTTL),
	}

	// todo: не создавать новую сессию на каждый запрос авторизации, если user_id, ip и user_agent совпадают

	session, err = s.sessionRepository.Create(ctx, session)
	if err != nil {
		s.logger.Error("create session", zap.Error(err))

		return nil, fmt.Errorf("create session: %w", err)
	}

	accessToken, err := s.SignAccessToken(user.ID, session.ID, s.config.Auth.AccessTokenTTL)
	if err != nil {
		s.logger.Error("sign access token", zap.Error(err))

		return nil, fmt.Errorf("sign access token: %w", err)
	}

	if err = s.userCacheRepository.Set(ctx, user, s.config.Auth.AccessTokenTTL); err != nil {
		s.logger.Error("set user to cache", zap.Error(err))
	}

	if err = s.sessionCacheRepository.Set(ctx, session, s.config.Auth.RefreshTokenTTL); err != nil {
		s.logger.Error("set session to cache", zap.Error(err))
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.Auth.AccessTokenTTL),
	}, nil
}

func (s *Service) Logout(ctx context.Context, userID int64, sessionID *string) error {
	if err := s.userCacheRepository.Delete(ctx, userID); err != nil {
		s.logger.Error("delete user cache", zap.Error(err))
	}

	query := &model.SessionQuery{
		Filter: model.SessionFilter{
			UserID: model.IDFilter{Eq: new(userID)},
		},
	}

	if sessionID != nil {
		query.Filter.ID = model.TextFilter{Eq: sessionID}
	}

	sessions, err := s.sessionRepository.List(ctx, query)
	if err != nil {
		s.logger.Error("get sessions list", zap.Error(err))

		return fmt.Errorf("get sessions list: %w", err)
	}

	for _, session := range sessions {
		if err = s.sessionCacheRepository.Delete(ctx, session.ID); err != nil {
			s.logger.Error("delete session cache", zap.Error(err))
		}

		if session.AuthType == model.AuthTypeVKOAuth {
			oauthAccessToken, err := s.encryptor.Decrypt(*session.OAuthAccessTokenEnc)
			if err != nil {
				s.logger.Error("decrypt oauth access token", zap.Error(err))

				return fmt.Errorf("decrypt oauth access token: %w", err)
			}

			if err = s.vkidClient.Logout(ctx, oauthAccessToken); err != nil {
				s.logger.Error("vkid logout", zap.Error(err))

				return fmt.Errorf("vkid logout: %w", err)
			}
		}
	}

	if err := s.sessionRepository.Delete(ctx, &query.Filter); err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil
		}

		s.logger.Error("delete session", zap.Error(err))

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
			return nil, model.ErrSessionNotFound
		}

		s.logger.Error("get session", zap.Error(err))

		return nil, fmt.Errorf("get session: %w", err)
	}

	// todo: решить, нужно ли сверять user_agent и ip, в случае, если они не совпадают, выбрасывать ошибку и удалять сессию

	if session.ExpiresAt.Before(time.Now()) {
		if err = s.sessionCacheRepository.Delete(ctx, session.ID); err != nil {
			s.logger.Error("delete session cache", zap.Error(err))
		}

		err = s.sessionRepository.Delete(ctx, &model.SessionFilter{
			ID:     model.TextFilter{Eq: new(session.ID)},
			UserID: model.IDFilter{Eq: new(session.UserID)},
		})
		if err != nil {
			if errors.Is(err, model.ErrSessionNotFound) {
				return nil, model.ErrSessionNotFound
			}

			s.logger.Error("delete session", zap.Error(err))

			return nil, fmt.Errorf("delete session: %w", err)
		}

		return nil, ErrTokenIsExpired
	}

	var (
		accessTokenExpiresIn  = s.config.Auth.AccessTokenTTL
		refreshTokenExpiresIn = s.config.Auth.RefreshTokenTTL
	)

	if session.AuthType == model.AuthTypeVKOAuth {
		oAuthRefreshToken, err := s.encryptor.Decrypt(*session.OAuthRefreshTokenEnc)
		if err != nil {
			s.logger.Error("decrypt oauth refresh token", zap.Error(err))

			return nil, fmt.Errorf("decrypt oauth refresh token: %w", err)
		}

		state, err := generateRandomString(stateBytes)
		if err != nil {
			s.logger.Error("generate state", zap.Error(err))

			return nil, fmt.Errorf("generate state: %w", err)
		}

		token, err := s.vkidClient.RefreshOAuthToken(ctx, &OAuthRefreshTokenParams{
			RefreshToken: oAuthRefreshToken,
			DeviceID:     *session.OAuthDeviceID,
			State:        state,
		})
		if err != nil {
			s.logger.Error("refresh oauth token", zap.Error(err))

			return nil, fmt.Errorf("refresh oauth toke: %w", err)
		}

		accessTokenExpiresIn = time.Duration(token.ExpiresIn) * time.Second
		refreshTokenExpiresIn = s.config.VKID.RefreshTokenTTL

		oAuthAccessToken, err := s.encryptor.Encrypt(token.AccessToken)
		if err != nil {
			s.logger.Error("encrypt oauth access token", zap.Error(err))

			return nil, fmt.Errorf("encrypt oauth access token: %w", err)
		}

		oAuthRefreshToken, err = s.encryptor.Encrypt(token.RefreshToken)
		if err != nil {
			s.logger.Error("encrypt oauth refresh token", zap.Error(err))

			return nil, fmt.Errorf("encrypt oauth refresh token: %w", err)
		}

		session.OAuthAccessTokenEnc = new(oAuthAccessToken)
		session.OAuthRefreshTokenEnc = new(oAuthRefreshToken)
	}

	refreshToken, err = generateRefreshToken()
	if err != nil {
		s.logger.Error("generate refresh token", zap.Error(err))

		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	session.ExpiresAt = time.Now().Add(refreshTokenExpiresIn)
	session.RefreshTokenHash = hashRefreshToken(refreshToken)

	accessToken, err := s.SignAccessToken(session.UserID, session.ID, accessTokenExpiresIn)
	if err != nil {
		s.logger.Error("sign access token", zap.Error(err))

		return nil, fmt.Errorf("sign access token: %w", err)
	}

	session, err = s.sessionRepository.Update(ctx, session)
	if err != nil {
		s.logger.Error("update session", zap.Error(err))

		return nil, fmt.Errorf("update session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(accessTokenExpiresIn),
	}, nil
}

func (s *Service) Authorize(
	ctx context.Context,
	token string,
) (*model.User, *model.Session, error) {
	claims, err := s.ParseAndValidateAccessToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("parse and validate auth token: %w", err)
	}

	var (
		user    *model.User
		session *model.Session
	)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		user, err = s.userCacheRepository.Get(ctx, claims.UserID)
		if err != nil {
			if !errors.Is(err, model.ErrUserNotFound) {
				s.logger.Error("get user from cache", zap.Error(err))
			}
		}

		if user != nil {
			return nil
		}

		user, err = s.userRepository.Get(ctx, &model.UserFilter{
			ID: model.IDFilter{Eq: new(claims.UserID)},
		})
		if err != nil {
			if errors.Is(err, model.ErrUserNotFound) {
				return model.ErrUserNotFound
			}

			return fmt.Errorf("get user: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		session, err = s.sessionCacheRepository.Get(ctx, claims.SessionID)
		if err != nil {
			if !errors.Is(err, model.ErrSessionNotFound) {
				s.logger.Error("get session from cache", zap.Error(err))
			}
		}

		if session != nil {
			return nil
		}

		session, err = s.sessionRepository.Get(ctx, &model.SessionFilter{
			ID: model.TextFilter{Eq: new(claims.SessionID)},
		})
		if err != nil {
			if errors.Is(err, model.ErrSessionNotFound) {
				return model.ErrSessionNotFound
			}

			return fmt.Errorf("get session: %w", err)
		}

		return nil
	})

	if err = group.Wait(); err != nil {
		s.logger.Error("authorize user", zap.Error(err))

		return nil, nil, fmt.Errorf("authorize user: %w", err)
	}

	return user, session, nil
}

func (s *Service) SignAccessToken(
	userID int64,
	sessionID string,
	expiresIn time.Duration,
) (string, error) {
	// todo: Возможно будет лучше перенести userID и sessionID в поля Subject и Issuer
	claims := AccessTokenClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
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
		return nil, ErrInvalidAuthToken
	}

	if claims, ok := t.Claims.(*AccessTokenClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, ErrInvalidAuthToken
}

func (s *Service) GenerateVKIDOAuthURL(ctx context.Context) (string, error) {
	codeVerifier, err := generateRandomString(codeVerifierBytes)
	if err != nil {
		s.logger.Error("generate code verifier", zap.Error(err))

		return "", fmt.Errorf("generate code verifier: %w", err)
	}

	state, err := generateRandomString(stateBytes)
	if err != nil {
		s.logger.Error("generate state", zap.Error(err))

		return "", fmt.Errorf("generate state: %w", err)
	}

	if err = s.codeVerifierRepository.Set(ctx, state, codeVerifier, codeVerifierTTL); err != nil {
		s.logger.Error("set code verifier", zap.Error(err))

		return "", fmt.Errorf("code verifier repository set: %w", err)
	}

	authURL, err := s.vkidClient.GenerateOAuthURL(codeChallengeS256(codeVerifier), state)
	if err != nil {
		s.logger.Error("generate oauth url", zap.Error(err))

		return "", fmt.Errorf("generate oauth url: %w", err)
	}

	return authURL, nil
}

func (s *Service) ExchangeVKIDOAuthToken(
	ctx context.Context,
	params *OAuthExchangeTokenParams,
) (*TokenPair, error) {
	codeVerifier, err := s.codeVerifierRepository.Pop(ctx, params.State)
	if err != nil {
		if errors.Is(err, model.ErrCodeVerifierNotFound) {
			return nil, model.ErrCodeVerifierNotFound
		}

		s.logger.Error("pop vkid code verifier", zap.Error(err))

		return nil, fmt.Errorf("pop vkid code verifier: %w", err)
	}

	params.CodeVerifier = codeVerifier

	oAuthToken, err := s.vkidClient.ExchangeOAuthToken(ctx, params)
	if err != nil {
		s.logger.Error("exchane oauth token", zap.Error(err))

		return nil, fmt.Errorf("exchane vkid oauth token: %w", err)
	}

	user, err := s.userRepository.Get(ctx, &model.UserFilter{
		VKUserID: model.IDFilter{Eq: new(oAuthToken.UserID)},
	})
	if err != nil {
		if !errors.Is(err, model.ErrUserNotFound) {
			s.logger.Error("get user by vk user id", zap.Error(err))

			return nil, fmt.Errorf("get user by vk user id: %w", err)
		}

		user, err = s.userRepository.Create(ctx, &model.User{
			VKUserID: new(oAuthToken.UserID),
			Role:     model.RoleUser,
		})
		if err != nil {
			s.logger.Error("create user with vk user id", zap.Error(err))

			return nil, fmt.Errorf("create user with vk user id: %w", err)
		}
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		s.logger.Error("generate refresh token", zap.Error(err))

		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	oAuthAccessTokenEnc, err := s.encryptor.Encrypt(oAuthToken.AccessToken)
	if err != nil {
		s.logger.Error("encrypt oauth access token", zap.Error(err))

		return nil, fmt.Errorf("encrypt oauth access token: %w", err)
	}

	oAuthRefreshTokenEnc, err := s.encryptor.Encrypt(oAuthToken.RefreshToken)
	if err != nil {
		s.logger.Error("encrypt oauth refresh token", zap.Error(err))

		return nil, fmt.Errorf("encrypt oauth refresh token: %w", err)
	}

	session := &model.Session{
		UserID:               user.ID,
		AuthType:             model.AuthTypeVKOAuth,
		IP:                   params.IP,
		UserAgent:            params.UserAgent,
		RefreshTokenHash:     hashRefreshToken(refreshToken),
		OAuthDeviceID:        new(params.DeviceID),
		OAuthAccessTokenEnc:  new(oAuthAccessTokenEnc),
		OAuthRefreshTokenEnc: new(oAuthRefreshTokenEnc),
		ExpiresAt:            time.Now().Add(s.config.VKID.RefreshTokenTTL),
	}

	// todo: не создавать новую сессию на каждый запрос авторизации, если user_id, ip и user_agent совпадают

	session, err = s.sessionRepository.Create(ctx, session)
	if err != nil {
		s.logger.Error("create session", zap.Error(err))

		return nil, fmt.Errorf("create session: %w", err)
	}

	accessTokenTTL := time.Duration(oAuthToken.ExpiresIn) * time.Second

	accessToken, err := s.SignAccessToken(user.ID, session.ID, accessTokenTTL)
	if err != nil {
		s.logger.Error("sign access token", zap.Error(err))

		return nil, fmt.Errorf("sign access token: %w", err)
	}

	if err = s.userCacheRepository.Set(ctx, user, accessTokenTTL); err != nil {
		s.logger.Error("set user to cache", zap.Error(err))
	}

	if err = s.sessionCacheRepository.Set(ctx, session, s.config.VKID.RefreshTokenTTL); err != nil {
		s.logger.Error("set session to cache", zap.Error(err))
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(accessTokenTTL),
	}, nil
}
