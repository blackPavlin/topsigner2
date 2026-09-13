package group

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/crypto"
	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger     *zap.Logger
	config     *config.Config
	encryptor  *crypto.Encryptor
	repository Repository
	vkClient   VKClient
}

func New(
	logger *zap.Logger,
	config *config.Config,
	encryptor *crypto.Encryptor,
	repository Repository,
	vkClient VKClient,
) *Service {
	return &Service{
		logger:     logger.Named("group-service"),
		config:     config,
		encryptor:  encryptor,
		repository: repository,
		vkClient:   vkClient,
	}
}

// todo: переделать возвращаемое значение на pagination.Result
func (s *Service) List(ctx context.Context, session *model.Session) ([]*Group, error) {
	token, err := s.encryptor.Decrypt(*session.OAuthAccessTokenEnc)
	if err != nil {
		s.logger.Error("decode access token", zap.Error(err))

		return nil, fmt.Errorf("decode access token: %w", err)
	}

	groups, err := s.vkClient.GetGroups(ctx, token)
	if err != nil {
		s.logger.Error("get vk groups", zap.Error(err))

		return nil, fmt.Errorf("get vk groups: %w", err)
	}

	return groups, nil
}
