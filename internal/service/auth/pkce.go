package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const (
	stateBytes        = 24
	codeVerifierBytes = 32
	codeVerifierTTL   = 10 * time.Minute
)

func codeChallengeS256(codeVerifier string) string {
	sum := sha256.Sum256([]byte(codeVerifier))

	return base64.RawURLEncoding.EncodeToString(sum[:])
}
