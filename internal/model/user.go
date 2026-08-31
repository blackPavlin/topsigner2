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
	ID           int64     `json:"id"`
	VKUserID     *int64    `json:"vk_user_id,omitempty"`
	Email        *string   `json:"email,omitempty"`
	PasswordHash *string   `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserFilter struct {
	ID       IDFilter
	VKUserID IDFilter
	Email    TextFilter
}
