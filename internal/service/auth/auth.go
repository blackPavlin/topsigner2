package auth

import (
	"errors"
	"time"
)

const (
	stateBytes      = 24
	codeVerifierTTL = 10 * time.Minute
)

var (
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrTokenIsExpired   = errors.New("token is expired")
	ErrInvalidPassword  = errors.New("invalid password")
)

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}
