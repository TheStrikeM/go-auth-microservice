package handlers

import (
	"microauth/internal/modules/auth/models/dto"
)

type IUserService interface {
	Register(userDto *dto.UserDTO) error
	SignIn(userDto *dto.UserDTO) (string, error)
}
