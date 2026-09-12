package group

import (
	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/config"
)

type Service struct {
	logger     *zap.Logger
	config     *config.Config
	repository Repository
	vkClient   VKClient
}

func New(
	logger *zap.Logger,
	config *config.Config,
	repository Repository,
	vkClient VKClient,
) *Service {
	return &Service{
		logger:     logger.Named("group-service"),
		config:     config,
		repository: repository,
		vkClient:   vkClient,
	}
}
