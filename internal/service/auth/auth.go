package auth

import "errors"

const stateBytes = 24

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
