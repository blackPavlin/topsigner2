package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

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

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash: %w", err)
	}

	return string(hash), nil
}

func ComparePasswordAndHash(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("compare password and hash: %w", err)
	}

	return nil
}
