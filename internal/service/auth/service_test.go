package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/crypto"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/auth"
	"github.com/bboykiv/topsigner/internal/service/auth/mock"
)

func TestService_Login_Success(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	password := "password123"

	passwordHash, err := crypto.GeneratePasswordHash(password)
	require.NoError(t, err)

	user := &model.User{
		ID:           1,
		Email:        new("test@email.com"),
		PasswordHash: new(passwordHash),
	}

	session := &model.Session{
		ID:     "session-id",
		UserID: user.ID,
	}

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(user, nil)

	sessionRepository.EXPECT().
		Create(t.Context(), gomock.Any()).
		Return(session, nil)

	userCacheRepository.EXPECT().
		Set(t.Context(), gomock.Any(), gomock.Any()).
		Return(nil)

	sessionCacheRepository.EXPECT().
		Set(t.Context(), gomock.Any(), gomock.Any()).
		Return(nil)

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	tokens, err := service.Login(t.Context(), &auth.LoginInput{
		Email:    *user.Email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)

	claims, err := service.ParseAndValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)
	require.Equal(t, user.ID, claims.UserID)
}

func TestService_Login_UserNotFound(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(nil, model.ErrUserNotFound)

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	tokens, err := service.Login(t.Context(), &auth.LoginInput{
		Email:    "test@email.com",
		Password: "password123",
	})
	require.ErrorIs(t, err, model.ErrUserNotFound)
	require.Nil(t, tokens)
}

func TestService_Login_InvalidPassword(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	user := &model.User{
		ID:           1,
		Email:        new("test@email.com"),
		PasswordHash: new("$2a$12$yyLAM0Pp2tZ/l3B4EK6IL.heTUqfnZnHiVq2lnCSoYMzSudD1cUX6"),
	}

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(user, nil)

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	tokens, err := service.Login(t.Context(), &auth.LoginInput{
		Email:    "test@email.com",
		Password: "invalid-password",
	})
	require.ErrorIs(t, err, auth.ErrInvalidPassword)
	require.Nil(t, tokens)
}

func TestService_Authorize_Success_EmptyCache(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	user := &model.User{
		ID: 1,
	}

	session := &model.Session{
		ID:       "session-id",
		UserID:   user.ID,
		AuthType: model.AuthTypePassword,
	}

	userCacheRepository.EXPECT().
		Get(gomock.Any(), user.ID).
		Return(nil, model.ErrUserNotFound)

	userRepository.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(user, nil)

	sessionCacheRepository.EXPECT().
		Get(gomock.Any(), session.ID).
		Return(nil, model.ErrSessionNotFound)

	sessionRepository.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(session, nil)

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	token, err := service.SignAccessToken(user.ID, session.ID, config.Auth.AccessTokenTTL)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authUser, authSession, err := service.Authorize(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, user.ID, authUser.ID)
	require.Equal(t, user.ID, authSession.UserID)
	require.Equal(t, model.AuthTypePassword, authSession.AuthType)
}

func TestService_Authorize_Success_NotEmptyCache(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	user := &model.User{
		ID: 1,
	}

	session := &model.Session{
		ID:       "session-id",
		UserID:   user.ID,
		AuthType: model.AuthTypePassword,
	}

	userCacheRepository.EXPECT().
		Get(gomock.Any(), user.ID).
		Return(user, nil)

	sessionCacheRepository.EXPECT().
		Get(gomock.Any(), session.ID).
		Return(session, nil)

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	token, err := service.SignAccessToken(user.ID, session.ID, config.Auth.AccessTokenTTL)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authUser, authSession, err := service.Authorize(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, user.ID, authUser.ID)
	require.Equal(t, user.ID, authSession.UserID)
	require.Equal(t, model.AuthTypePassword, authSession.AuthType)
}

func TestService_Authorize_InvatidToken(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	user, session, err := service.Authorize(t.Context(), "invalid token")
	require.Nil(t, user)
	require.Nil(t, session)
	require.ErrorIs(t, err, auth.ErrInvalidAuthToken)
}

func TestService_Authorize_UserNotFound(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		sessionCacheRepository = mock.NewMockSessionCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			RefreshTokenTTL: 5 * time.Minute,
			AccessTokenTTL:  1 * time.Minute,
			SigningKey:      "secret-key",
			EncryptionKey:   "8f3kL9mN2pQrS7tVxYzA1bC4dE5fG6hI7jK8lM9nO0p=",
		},
	}

	var (
		userID    int64 = 1
		sessionID       = "session-id"
	)

	userCacheRepository.EXPECT().
		Get(gomock.Any(), userID).
		Return(nil, model.ErrUserNotFound)

	userRepository.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil, model.ErrUserNotFound)

	sessionCacheRepository.EXPECT().
		Get(gomock.Any(), sessionID).
		Return(&model.Session{}, nil)

	encryptor, err := crypto.NewEncryptor(config.Auth.EncryptionKey)
	require.NoError(t, err)

	service := auth.New(
		zap.NewNop(),
		config,
		encryptor,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		sessionCacheRepository,
		codeVerifierRepository,
	)

	token, err := service.SignAccessToken(userID, sessionID, config.Auth.AccessTokenTTL)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	user, session, err := service.Authorize(t.Context(), token)
	require.ErrorIs(t, err, model.ErrUserNotFound)
	require.Nil(t, user)
	require.Nil(t, session)
}
