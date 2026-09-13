package crypto_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/internal/crypto"
)

func generateBase64Key(t *testing.T, keyLen int) string {
	t.Helper()

	key := make([]byte, keyLen)
	_, err := rand.Read(key)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(key)
}

func TestNewEncryptor_Success(t *testing.T) {
	key := generateBase64Key(t, 32)

	encryptor, err := crypto.NewEncryptor(key)
	require.NoError(t, err)
	require.NotNil(t, encryptor)
}

func TestNewEncryptor_Failed(t *testing.T) {
	type testCase struct {
		name string
		key  string
	}

	testCases := []testCase{
		{
			name: "invalid base64 returns error",
			key:  "not-valid-base64-!!!",
		},
		{
			name: "key shorter than 32 bytes",
			key:  generateBase64Key(t, 16),
		},
		{
			name: "key longer than 32 bytes",
			key:  generateBase64Key(t, 33),
		},
		{
			name: "empty key",
			key:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptor, err := crypto.NewEncryptor(tc.key)

			require.Error(t, err)
			require.Nil(t, encryptor)
		})
	}
}

func TestEncryptor_EncryptDecrypt_RoundTrip(t *testing.T) {
	encryptor, err := crypto.NewEncryptor(generateBase64Key(t, 32))
	require.NoError(t, err)

	type testCase struct {
		name      string
		plaintext string
	}

	testCases := []testCase{
		{
			name:      "simple ascii string",
			plaintext: "hello world",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode string",
			plaintext: "Hallo, world! 🔐",
		},
		{
			name:      "long string",
			plaintext: string(make([]byte, 10000)),
		},
		{
			name:      "oauth-like token",
			plaintext: "vk1.a.AbCdEf1234567890_-ExampleToken",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext, err := encryptor.Encrypt(tc.plaintext)
			require.NoError(t, err)
			require.NotEmpty(t, ciphertext)

			decrypted, err := encryptor.Decrypt(ciphertext)
			require.NoError(t, err)
			require.Equal(t, tc.plaintext, decrypted)
		})
	}
}
