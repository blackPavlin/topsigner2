package model

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("user must have email or vk id")
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID           int64
	VKUserID     *int64
	Email        *string
	PasswordHash *string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserFilter struct {
	ID       IDFilter
	VKUserID IDFilter
	Email    TextFilter
}
