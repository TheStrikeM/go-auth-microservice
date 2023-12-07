package handlers

import (
	"microauth/internal/domain/models"
	"microauth/internal/modules/user/models/dto"
)

type IUserService interface {
	UserById(id string) (*models.User, error)
	UserByUsername(username string) (*models.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, dto *dto.UpdateUserDto) (*models.User, error)
}
