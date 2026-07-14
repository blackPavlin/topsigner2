package model

import (
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

type Session struct {
	ID               string
	UserID           int64
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SessionFilter struct {
	ID               IDFilter
	UserID           IDFilter
	RefreshTokenHash TextFilter
}
