package group

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/bboykiv/topsigner/internal/adapter/crypto"
	"github.com/bboykiv/topsigner/internal/config"
	"github.com/bboykiv/topsigner/internal/model"
)

type Service struct {
	logger            *zap.Logger
	config            *config.Config
	encryptor         *crypto.Encryptor
	groupRepository   Repository
	sessionRepository SessionRepository
	stateRepository   StateRepository
	vkClient          VKClient
}

func New(
	logger *zap.Logger,
	config *config.Config,
	encryptor *crypto.Encryptor,
	groupRepository Repository,
	sessionRepository SessionRepository,
	stateRepository StateRepository,
	vkClient VKClient,
) *Service {
	return &Service{
		logger:            logger.Named("group-service"),
		config:            config,
		encryptor:         encryptor,
		groupRepository:   groupRepository,
		sessionRepository: sessionRepository,
		stateRepository:   stateRepository,
		vkClient:          vkClient,
	}
}

// todo: переделать метод на полноценную пагинацию
func (s *Service) List(
	ctx context.Context,
	sessionID string,
	query *model.GroupQuery,
) (*model.List[*Group], error) {
	session, err := s.sessionRepository.Get(ctx, &model.SessionFilter{
		ID: model.TextFilter{Eq: new(sessionID)},
	})
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, model.ErrSessionNotFound
		}

		s.logger.Error("get user session", zap.Error(err))

		return nil, fmt.Errorf("get user session: %w", err)
	}

	// todo: добавить проверку, что OAuthAccessTokenEnc != nil и AuthType = VK_OAUTH

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

	connectedGroups, err := s.groupRepository.List(ctx, &model.GroupQuery{
		Filter: model.GroupFilter{
			UserID: model.IDFilter{Eq: new(session.UserID)},
		},
		Pagination: model.Pagination{
			Limit: model.DefaultPaginationLimit,
		},
	})
	if err != nil {
		s.logger.Error("get groups", zap.Error(err))

		return nil, fmt.Errorf("get groups: %w", err)
	}

	connectedExternalIDs := make(map[int64]struct{}, len(connectedGroups))

	for _, connectedGroup := range connectedGroups {
		connectedExternalIDs[connectedGroup.ExternalID] = struct{}{}
	}

	for _, group := range groups {
		if _, ok := connectedExternalIDs[group.ID]; ok {
			group.IsConnected = true
		}
	}

	result := &model.List[*Group]{
		Items: groups,
	}

	return result, nil
}

func (s *Service) Connect(ctx context.Context, code, state string) (*model.Group, error) {
	userID, err := s.stateRepository.Pop(ctx, state)
	if err != nil {
		if errors.Is(err, model.ErrGroupStateNotFound) {
			return nil, model.ErrGroupStateNotFound
		}

		s.logger.Error("pop group state", zap.Error(err))

		return nil, fmt.Errorf("pop group state: %w", err)
	}

	groupID, token, err := s.vkClient.ExchangeGroupCode(ctx, code)
	if err != nil {
		s.logger.Error("exchange group code", zap.Error(err))

		return nil, fmt.Errorf("exchange group token: %w", err)
	}

	accessTokenEnc, err := s.encryptor.Encrypt(token)
	if err != nil {
		s.logger.Error("encrypt group access token", zap.Error(err))

		return nil, fmt.Errorf("encrypt group access token: %w", err)
	}

	group, err := s.groupRepository.Create(ctx, &model.Group{
		UserID:         userID,
		ExternalID:     groupID,
		AccessTokenEnc: accessTokenEnc,
	})
	if err != nil {
		s.logger.Error("create group", zap.Error(err))

		return nil, fmt.Errorf("create group: %w", err)
	}

	return group, nil
}

func (s *Service) GenerateConnectionURL(ctx context.Context, userID, groupID int64) (string, error) {
	state, err := crypto.GenerateRandomString(stateSize)
	if err != nil {
		s.logger.Error("generate state", zap.Error(err))

		return "", fmt.Errorf("generate state: %w", err)
	}

	if err = s.stateRepository.Set(ctx, state, userID, s.config.VK.OAuthStateTTL); err != nil {
		s.logger.Error("set group state", zap.Error(err))

		return "", fmt.Errorf("set group state: %w", err)
	}

	connectionURL, err := s.vkClient.GenerateConnectGroupURL(groupID, state)
	if err != nil {
		s.logger.Error("generate connection url", zap.Error(err))

		return "", fmt.Errorf("generate connection url: %w", err)
	}

	return connectionURL, nil
}
