package keyvalue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bboykiv/topsigner/internal/model"
)

const sessionCachePrefix = "session_cache:"

type SessionCacheRepository struct {
	client *redis.Client
}

func NewSessionCacheRepository(client *redis.Client) *SessionCacheRepository {
	return &SessionCacheRepository{client: client}
}

func (r *SessionCacheRepository) Get(
	ctx context.Context,
	sessionID string,
) (*model.Session, error) {
	result, err := r.client.Get(ctx, getSessionCacheKey(sessionID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, model.ErrSessionNotFound
		}

		return nil, fmt.Errorf("get session from cache: %w", err)
	}

	var session model.Session

	if err = json.Unmarshal([]byte(result), &session); err != nil {
		return nil, fmt.Errorf("unmarshal get session cache result: %w", err)
	}

	return &session, nil
}

func (r *SessionCacheRepository) Set(
	ctx context.Context,
	session *model.Session,
	ttl time.Duration,
) error {
	result, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal cache session: %w", err)
	}

	err = r.client.Set(ctx, getSessionCacheKey(session.ID), result, ttl).Err()
	if err != nil {
		return fmt.Errorf("set session to cache: %w", err)
	}

	return nil
}

func (r *SessionCacheRepository) Delete(ctx context.Context, sessionID string) error {
	if err := r.client.Del(ctx, getSessionCacheKey(sessionID)).Err(); err != nil {
		return fmt.Errorf("delete session cache: %w", err)
	}

	return nil
}

func getSessionCacheKey(sessionID string) string {
	return sessionCachePrefix + sessionID
}
