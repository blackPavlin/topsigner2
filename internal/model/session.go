package model

import (
	"errors"
	"time"
)

type AuthType string

const (
	AuthTypePassword AuthType = "PASSWORD"
	AuthTypeVKOAuth  AuthType = "VK_OAUTH"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrCodeVerifierNotFound = errors.New("code verifier not found")
)

type Session struct {
	ID                   string
	UserID               int64
	AuthType             AuthType
	IP                   string
	UserAgent            string
	RefreshTokenHash     string
	OAuthDeviceID        *string
	OAuthAccessTokenEnc  *string
	OAuthRefreshTokenEnc *string
	ExpiresAt            time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type SessionFilter struct {
	ID               TextFilter
	UserID           IDFilter
	AuthType         TextFilter
	RefreshTokenHash TextFilter
}

type SessionQuery struct {
	Filter SessionFilter
}
