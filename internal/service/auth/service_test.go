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
	"github.com/bboykiv/topsigner/internal/service/auth/mocks"
)

func TestService_Login_Success(t *testing.T) {
	var (
		ctrl              = gomock.NewController(t)
		userRepository    = mocks.NewMockUserRepository(ctrl)
		sessionRepository = mocks.NewMockSessionRepository(ctrl)
		vkidClient        = mocks.NewMockVKIDClient(ctrl)
	)

	config := &config.Config{
		Auth: config.AuthConfig{
			AccessTokenTTL: 15 * time.Minute,
			SigningKey:     "secret-key",
		},
	}

	password := "password123"

	passwordHash, err := auth.GeneratePasswordHash(password)
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

	svc := auth.New(zap.NewNop(), config, vkidClient, userRepository, sessionRepository)

	tokens, err := svc.Login(t.Context(), &auth.LoginInput{
		Email:    user.Email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)

	claims, err := svc.ParseAndValidateAccessToken(tokens.AccessToken)
	require.NoError(t, err)
	require.Equal(t, user.ID, claims.UserID)
}

func TestService_Login_UserNotFound(t *testing.T) {
	var (
		ctrl              = gomock.NewController(t)
		userRepository    = mocks.NewMockUserRepository(ctrl)
		sessionRepository = mocks.NewMockSessionRepository(ctrl)
		vkidClient        = mocks.NewMockVKIDClient(ctrl)
	)

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(nil, model.ErrUserNotFound)

	svc := auth.New(zap.NewNop(), &config.Config{}, vkidClient, userRepository, sessionRepository)

	tokens, err := svc.Login(t.Context(), &auth.LoginInput{
		Email:    "test@email.com",
		Password: "password123",
	})
	require.ErrorIs(t, err, model.ErrUserNotFound)
	require.Nil(t, tokens)
}

func TestService_Login_InvalidPassword(t *testing.T) {
	var (
		ctrl              = gomock.NewController(t)
		userRepository    = mocks.NewMockUserRepository(ctrl)
		sessionRepository = mocks.NewMockSessionRepository(ctrl)
		vkidClient        = mocks.NewMockVKIDClient(ctrl)
	)

	user := &model.User{
		ID:           1,
		Email:        "test@email.com",
		PasswordHash: "$2a$12$yyLAM0Pp2tZ/l3B4EK6IL.heTUqfnZnHiVq2lnCSoYMzSudD1cUX6",
	}

	userRepository.EXPECT().
		Get(t.Context(), gomock.Any()).
		Return(user, nil)

	svc := auth.New(zap.NewNop(), &config.Config{}, vkidClient, userRepository, sessionRepository)

	tokens, err := svc.Login(t.Context(), &auth.LoginInput{
		Email:    "test@email.com",
		Password: "invalid-password",
	})
	require.ErrorIs(t, err, auth.ErrInvalidPassword)
	require.Nil(t, tokens)
}
