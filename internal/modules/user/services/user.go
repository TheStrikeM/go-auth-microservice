package services

import (
	"context"
	"log/slog"
	"microauth/internal/domain/models"
	"microauth/internal/modules/user/models/dto"
	"time"
)

type IUserRepository interface {
	UserById(ctx context.Context, id string) (*models.User, error)
	UserByUsername(ctx context.Context, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, dto *dto.UpdateUserDto) (*models.User, error)
}

type UserService struct {
	repo IUserRepository
	log  *slog.Logger
}

func New(repo IUserRepository, log *slog.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (us *UserService) UserById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := us.repo.UserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) UserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := us.repo.UserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := us.repo.DeleteUser(ctx, id); err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdateUser(id string, dto *dto.UpdateUserDto) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := us.repo.UpdateUser(ctx, id, dto)
	if err != nil {
		return nil, err
	}

	return user, nil
}
