package model

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID        int64
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserFilter struct {
	ID IDFilter
}
