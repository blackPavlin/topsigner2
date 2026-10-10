package model

import (
	"errors"
	"time"
	"uuid"
)

var (
	ErrVKSessionNotFound    = errors.New("vk session not found")
	ErrCodeVerifierNotFound = errors.New("code verifier not found")
)

type VKSession struct {
	ID              uuid.UUID
	SessionID       uuid.UUID
	DeviceID        string
	AccessTokenEnc  string
	RefreshTokenEnc string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type VKSessionFilter struct {
	ID        UUIDFileter
	SessionID UUIDFileter
}
