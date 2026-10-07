package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomString(size int) (string, error) {
	buffer := make([]byte, size)

	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate random string: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
