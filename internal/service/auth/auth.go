package auth

import (
	"errors"
)

var (
	ErrInvalidAuthToken          = errors.New("invalid auth token")
	ErrTokenIsExpired            = errors.New("token is expired")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrPasswordLoginNotAvailable = errors.New("password login not available")
)

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}
