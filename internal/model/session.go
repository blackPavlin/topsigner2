package model

import (
	"errors"
	"time"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrCodeVerifierNotFound = errors.New("code verifier not found")
)

type Session struct {
	ID               string
	UserID           int64
	IP               string
	UserAgent        string
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SessionFilter struct {
	ID               TextFilter
	UserID           IDFilter
	IP               TextFilter
	RefreshTokenHash TextFilter
}
