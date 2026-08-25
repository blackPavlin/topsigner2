package keyvalue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bboykiv/topsigner/internal/model"
)

const userCachePrefix = "user_cache:"

type UserCacheRepository struct {
	client *redis.Client
}

func NewUserCacheRepository(client *redis.Client) *UserCacheRepository {
	return &UserCacheRepository{client: client}
}

func (r *UserCacheRepository) Get(ctx context.Context, userID int64) (*model.User, error) {
	result, err := r.client.Get(ctx, getUserCacheKey(userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, model.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user from cache: %w", err)
	}

	var user model.User

	if err = json.Unmarshal([]byte(result), &user); err != nil {
		return nil, fmt.Errorf("unmarshal get user cache result: %w", err)
	}

	return &user, nil
}

func (r *UserCacheRepository) Set(ctx context.Context, user *model.User, ttl time.Duration) error {
	result, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal cache user: %w", err)
	}

	err = r.client.Set(ctx, getUserCacheKey(user.ID), result, ttl).Err()
	if err != nil {
		return fmt.Errorf("set user to cache: %w", err)
	}

	return nil
}

func (r *UserCacheRepository) Delete(ctx context.Context, userID int64) error {
	if err := r.client.Del(ctx, getUserCacheKey(userID)).Err(); err != nil {
		return fmt.Errorf("delete user cache: %w", err)
	}

	return nil
}

func getUserCacheKey(userID int64) string {
	return userCachePrefix + strconv.FormatInt(userID, 10)
}
