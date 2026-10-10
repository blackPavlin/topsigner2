package auth

import (
	"context"
	"uuid"

	"github.com/bboykiv/topsigner/internal/model"
)

type SessionManager struct {
	sessionRepository SessionRepository
}

func NewSessionManager(
	sessionRepository SessionRepository,
) *SessionManager {
	return &SessionManager{
		sessionRepository: sessionRepository,
	}
}

// Выдать сессию после успешной аутентификации (любым способом)
func (m *SessionManager) Start(ctx context.Context, user *model.User, p StartParams) (*Token, *model.Session, error)

// Обменять refresh-токен на новую пару с ротацией
func (m *SessionManager) Refresh(ctx context.Context, refreshToken, deviceID string) (*Token, error)

// Logout одной сессии / всех сессий пользователя
func (m *SessionManager) Revoke(ctx context.Context, sessionID uuid.UUID) error
func (m *SessionManager) RevokeAll(ctx context.Context, userID int64) error
