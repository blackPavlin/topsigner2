package model

import (
	"errors"
	"time"
)

var (
	ErrUnsupportedImageFormat = errors.New("unsupported image format")
	ErrImageAlreadyExists     = errors.New("image already exists")
)

type Image struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
