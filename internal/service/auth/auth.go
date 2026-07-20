package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const refreshTokenBytes = 32

var (
	ErrInvalidAuthToken = errors.New("invalid auth token")
	ErrTokenIsExpired   = errors.New("token is expired")
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserID    int64  `json:"user_id"`
	SessionID string `json:"sid"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type LoginInput struct {
	Email string
	IP    string
}

type LogoutInput struct {}
