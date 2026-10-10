package model

import "errors"

var (
	ErrPasswordLoginNotAvailable = errors.New("password login not available")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrInvalidAuthCode           = errors.New("invalid auth code")
	ErrUpstreamUnavailable       = errors.New("upstream unavailable")
)
