package user_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
	"github.com/bboykiv/topsigner/internal/service/user"
	"github.com/bboykiv/topsigner/internal/service/user/mock"
)

func TestService_Create_Success(t *testing.T) {
	var (
		logger     = zap.NewNop()
		ctrl       = gomock.NewController(t)
		repository = mock.NewMockRepository(ctrl)
	)

	input := &user.CreateUserInput{
		Email:    "test@email.com",
		Password: "password",
		Role:     model.RoleUser,
	}

	repository.EXPECT().
		Create(t.Context(), gomock.Any()).
		Return(&model.User{
			ID:    1,
			Email: new(input.Email),
			Role:  input.Role,
		}, nil)

	service := user.New(logger, &config.Config{}, repository)

	user, err := service.Create(t.Context(), input)
	require.NoError(t, err)
	require.NotNil(t, user)
}

func TestService_CreateDefault_Success(t *testing.T) {
	var (
		logger     = zap.NewNop()
		ctrl       = gomock.NewController(t)
		repository = mock.NewMockRepository(ctrl)
	)

	userConfig := config.UserConfig{
		Default: config.DefaultUserConfig{
			Email:    "admin@email.com",
			Password: "password",
			Role:     model.RoleAdmin,
		},
	}

	repository.EXPECT().
		Create(t.Context(), gomock.Any()).
		Return(&model.User{
			ID:    1,
			Email: new(userConfig.Default.Email),
			Role:  userConfig.Default.Role,
		}, nil)

	service := user.New(logger, &config.Config{User: userConfig}, repository)

	err := service.CreateDefault(t.Context())
	require.NoError(t, err)
}

func TestService_CreateDefault_AlreadyExists(t *testing.T) {
	var (
		logger     = zap.NewNop()
		ctrl       = gomock.NewController(t)
		repository = mock.NewMockRepository(ctrl)
	)

	userConfig := config.UserConfig{
		Default: config.DefaultUserConfig{
			Email:    "admin@email.com",
			Password: "password",
			Role:     model.RoleAdmin,
		},
	}

	repository.EXPECT().
		Create(t.Context(), gomock.Any()).
		Return(nil, model.ErrUserAlreadyExists)

	service := user.New(logger, &config.Config{User: userConfig}, repository)

	err := service.CreateDefault(t.Context())
	require.NoError(t, err)
}
