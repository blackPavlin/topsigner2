package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bboykiv/topsigner/internal/model"
)

const codeVerifierPrefix = "code_verifier:"

type CodeVerifierStore struct {
	client *redis.Client
}

func NewCodeVerifierStore(client *redis.Client) *CodeVerifierStore {
	return &CodeVerifierStore{client: client}
}

func getCodeVerifierKey(state string) string {
	return codeVerifierPrefix + state
}

func (r *CodeVerifierStore) Set(
	ctx context.Context,
	state, verifier string,
	ttl time.Duration,
) error {
	if err := r.client.Set(ctx, getCodeVerifierKey(state), verifier, ttl).Err(); err != nil {
		return fmt.Errorf("set code verifier: %w", err)
	}

	return nil
}

func (r *CodeVerifierStore) Pop(ctx context.Context, state string) (string, error) {
	result, err := r.client.GetDel(ctx, getCodeVerifierKey(state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.ErrCodeVerifierNotFound
		}

		return "", fmt.Errorf("pop code verifier: %w", err)
	}

	return result, nil
}
