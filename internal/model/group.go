package model

import (
	"errors"
	"time"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupAlreadyExists = errors.New("group already exists")
)

type Group struct {
	ID         int64
	UserID     int64
	ExternalID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type GroupFilter struct {
	ID     IDFilter
	UserID IDFilter
}

type GroupQuery struct {
	Filter     GroupFilter
	Pagination Pagination
}
