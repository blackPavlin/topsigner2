package model

import "errors"

var (
	ErrInvalidAuthCode     = errors.New("invalid auth code")
	ErrUpstreamUnavailable = errors.New("upstream unavailable")
)
