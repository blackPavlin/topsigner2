package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bboykiv/topsigner/internal/model"
)

func TestEncodeCursor(t *testing.T) {
	type testCase struct {
		Name      string
		ID        int64
		CreatedAt time.Time
		Result    string
	}

	tests := []testCase{
		{
			Name:      "correct encoding",
			ID:        1,
			CreatedAt: time.Date(2026, 7, 3, 10, 30, 0, 0, time.UTC),
			Result:    "MV8xNzgzMDc0NjAwMDAwMDAwMDAw",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cursor := model.EncodeCursor(tc.ID, tc.CreatedAt)

			require.Equal(t, tc.Result, cursor)
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	type testCase struct {
		Name   string
		Cursor string
		Result *model.Cursor
		Error  error
	}

	tests := []testCase{
		{
			Name:   "correct decoding",
			Cursor: "MV8xNzgzMDc0NjAwMDAwMDAwMDAw",
			Result: &model.Cursor{
				ID:        1,
				CreatedAt: time.Date(2026, 7, 3, 10, 30, 0, 0, time.UTC),
			},
			Error: nil,
		},
		{
			Name:   "empty cursor string",
			Cursor: "",
			Result: nil,
			Error:  model.ErrInvalidCursor,
		},
		{
			Name:   "invalid cursor string",
			Cursor: "NA==",
			Result: nil,
			Error:  model.ErrInvalidCursor,
		},
		{
			Name:   "invalid cursor id",
			Cursor: "ZXJyb3JfMTc4MjQ5ODA5MDcyMDk5NDAwMA==",
			Result: nil,
			Error:  model.ErrInvalidCursor,
		},
		{
			Name:   "invalid cursor date",
			Cursor: "MV9lcnJvcjE3ODI0OTgwOTA3MjA5OTQwMDA=",
			Result: nil,
			Error:  model.ErrInvalidCursor,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cursor, err := model.DecodeCursor(tc.Cursor)
			if err != nil {
				require.ErrorIs(t, err, tc.Error)

				return
			}

			require.NotNil(t, cursor)
			require.Equal(t, cursor.ID, tc.Result.ID)
			require.Equal(t, cursor.CreatedAt, tc.Result.CreatedAt)
		})
	}
}
