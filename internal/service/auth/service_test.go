package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
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
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SigningKey:     "secret-key",
		},
	}

	password := "password123"

	passwordHash, err := model.GeneratePasswordHash(password)
	require.NoError(t, err)

	user := &model.User{
		ID:           1,
		Email:        "test@email.com",
		PasswordHash: passwordHash,
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

	service := auth.New(
		zap.NewNop(),
		config,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		codeVerifierRepository,
	)

	tokens, err := service.Login(t.Context(), &auth.LoginInput{
		Email:    user.Email,
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
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(nil, model.ErrUserNotFound)

	service := auth.New(
		zap.NewNop(),
		&config.Config{},
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
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
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	user := &model.User{
		ID:           1,
		Email:        "test@email.com",
		PasswordHash: "$2a$12$yyLAM0Pp2tZ/l3B4EK6IL.heTUqfnZnHiVq2lnCSoYMzSudD1cUX6",
	}

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(user, nil)

	service := auth.New(
		zap.NewNop(),
		&config.Config{},
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
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
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SigningKey:     "secret-key",
		},
	}

	user := &model.User{
		ID: 1,
	}

	session := &model.Session{
		ID: "session-id",
	}

	userCacheRepository.EXPECT().
		Get(t.Context(), user.ID).
		Return(nil, model.ErrUserNotFound)

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(user, nil)

	userCacheRepository.EXPECT().
		Set(t.Context(), gomock.Any(), gomock.Any()).
		Return(nil)

	service := auth.New(
		zap.NewNop(),
		config,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		codeVerifierRepository,
	)

	token, err := service.SignAccessToken(user.ID, session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizedUser, err := service.Authorize(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, user.ID, authorizedUser.ID)
}

func TestService_Authorize_Success_NotEmptyCache(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SigningKey:     "secret-key",
		},
	}

	user := &model.User{
		ID: 1,
	}

	session := &model.Session{
		ID: "session-id",
	}

	userCacheRepository.EXPECT().
		Get(t.Context(), user.ID).
		Return(user, nil)

	service := auth.New(
		zap.NewNop(),
		config,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		codeVerifierRepository,
	)

	token, err := service.SignAccessToken(user.ID, session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizedUser, err := service.Authorize(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, user.ID, authorizedUser.ID)
}

func TestService_Authorize_InvatidToken(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SigningKey:     "secret-key",
		},
	}

	service := auth.New(
		zap.NewNop(),
		config,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		codeVerifierRepository,
	)

	user, err := service.Authorize(t.Context(), "invalid token")
	require.Nil(t, user)
	require.ErrorIs(t, err, auth.ErrInvalidAuthToken)
}

func TestService_Authorize_UserNotFound(t *testing.T) {
	var (
		ctrl                   = gomock.NewController(t)
		userRepository         = mock.NewMockUserRepository(ctrl)
		sessionRepository      = mock.NewMockSessionRepository(ctrl)
		codeVerifierRepository = mock.NewMockCodeVerifierRepository(ctrl)
		userCacheRepository    = mock.NewMockUserCacheRepository(ctrl)
		vkidClient             = mock.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SigningKey:     "secret-key",
		},
	}

	var (
		userID    int64 = 1
		sessionID       = "session-id"
	)

	userCacheRepository.EXPECT().
		Get(t.Context(), userID).
		Return(nil, model.ErrUserNotFound)

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(nil, model.ErrUserNotFound)

	service := auth.New(
		zap.NewNop(),
		config,
		vkidClient,
		userRepository,
		sessionRepository,
		userCacheRepository,
		codeVerifierRepository,
	)

	token, err := service.SignAccessToken(userID, sessionID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	user, err := service.Authorize(t.Context(), token)
	require.ErrorIs(t, err, model.ErrUserNotFound)
	require.Nil(t, user)
}
