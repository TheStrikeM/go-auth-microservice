package services

import (
	"context"
	"log/slog"
	"microauth/internal/modules/auth/models/dto"
	"microauth/internal/modules/auth/models/models"
)

type IUserRepo interface {
	AddUser(ctx context.Context, userDto *dto.UserDTO) error
	User(ctx context.Context, userDto *dto.UserDTO) (*models.User, error)
}

type UserService struct {
	userRepo IUserRepo
	log      *slog.Logger
}

func New(userRepo IUserRepo, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		log:      logger,
	}
}

func (us *UserService) SignIn(userDto *dto.UserDTO) (string, error) {
	//validate
	//check is true on database
	//generate jwt token
	//return token
}

func (us *UserService) Register(userDto *dto.UserDTO) error {
	//validate
	//check unique
	//add to database
}
