package model

import (
	"errors"
	"net/netip"
	"time"
	"uuid"
)

type SessionAuthType string

const (
	SessionAuthTypePassword SessionAuthType = "PASSWORD"
	SessionAuthTypeVKOAuth  SessionAuthType = "VK_OAUTH"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	ID               uuid.UUID
	UserID           int64
	AuthType         SessionAuthType
	DeviceID         string
	IP               netip.Addr
	UserAgent        string
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SessionFilter struct {
	ID               UUIDFileter
	UserID           IDFilter
	AuthType         TextFilter
	RefreshTokenHash TextFilter
}

type SessionQuery struct {
	Filter SessionFilter
}
