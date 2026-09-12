package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	stateBytes        = 24
	codeVerifierBytes = 32
	codeVerifierTTL   = 10 * time.Minute
)

func generateRandomString(size int) (string, error) {
	buffer := make([]byte, size)

	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate random string: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func codeChallengeS256(codeVerifier string) string {
	sum := sha256.Sum256([]byte(codeVerifier))

	return base64.RawURLEncoding.EncodeToString(sum[:])
}
