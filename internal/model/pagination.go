package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DefaultPaginationLimit = 50

var ErrInvalidCursor = errors.New("invalid cursor")

type Pagination struct {
	Cursor *Cursor
	Limit  int
}

type Cursor struct {
	ID        int64
	CreatedAt time.Time
}

type List[T any] struct {
	Items      []T
	NextCursor *string
	HasNext    bool
}

func EncodeCursor(id int64, createdAt time.Time) string {
	payload := fmt.Sprintf("%d_%d", id, createdAt.UnixNano())

	return base64.URLEncoding.EncodeToString([]byte(payload))
}

func DecodeCursor(cursor string) (*Cursor, error) {
	decoded, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}

	parts := strings.Split(string(decoded), "_")
	if len(parts) != 2 {
		return nil, ErrInvalidCursor
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	createdAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidCursor
	}

	return &Cursor{
		ID:        id,
		CreatedAt: time.Unix(0, createdAt).UTC(),
	}, nil
}
