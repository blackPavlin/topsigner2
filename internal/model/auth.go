package model

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidAuthToken = errors.New("invalid auth token")

type AuthClaims struct {
	jwt.RegisteredClaims
	UserID int64
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
}
