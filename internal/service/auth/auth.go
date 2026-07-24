package auth

import "errors"

const stateBytes = 24

var (
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrTokenIsExpired   = errors.New("token is expired")
)

type LoginInput struct {
	Email string
	IP    string
}
