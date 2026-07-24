package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const refreshTokenBytes = 32

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserID    int64  `json:"user_id"`
	SessionID string `json:"sid"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func generateRefreshToken() (string, error) {
	token, err := generateRandomString(refreshTokenBytes)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return token, nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}
