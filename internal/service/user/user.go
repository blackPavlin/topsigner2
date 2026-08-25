package user

import "github.com/bboykiv/topsigner/internal/model"

type CreateUserInput struct {
	Email    string
	Password string
	Role     model.Role
}
