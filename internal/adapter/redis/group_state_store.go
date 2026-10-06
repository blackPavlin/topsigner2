package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/bboykiv/topsigner/internal/model"
)

const groupStatePrefix = "group_state:"

type GroupStateStore struct {
	client *redis.Client
}

func NewGroupStateStore(client *redis.Client) *GroupStateStore {
	return &GroupStateStore{client: client}
}

func getGroupStateKey(state string) string {
	return groupStatePrefix + state
}

func (r *GroupStateStore) Set(
	ctx context.Context,
	state string,
	userID int64,
	ttl time.Duration,
) error {
	err := r.client.Set(ctx, getGroupStateKey(state), strconv.FormatInt(userID, 10), ttl).Err()
	if err != nil {
		return fmt.Errorf("set group state: %w", err)
	}

	return nil
}

func (r *GroupStateStore) Pop(ctx context.Context, state string) (int64, error) {
	result, err := r.client.GetDel(ctx, getGroupStateKey(state)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, model.ErrGroupStateNotFound
		}

		return 0, fmt.Errorf("getdel group state: %w", err)
	}

	userID, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse user id: %w", err)
	}

	return userID, nil
}
