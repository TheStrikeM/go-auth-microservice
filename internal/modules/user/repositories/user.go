package repositories

import (
	"context"
	"microauth/internal/domain/models"
	"microauth/internal/modules/user/models/dto"
	"microauth/internal/storage/clients/postgresql"
)

type IUserRepository interface {
	UserById(ctx context.Context, id string) (*models.User, error)
	UserByUsername(ctx context.Context, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, dto *dto.UpdateUserDto) (*models.User, error)
}

type UserRepository struct {
	client *postgresql.PSQLClient
}

func New(client *postgresql.PSQLClient) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (ur *UserRepository) UserById(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT *
		FROM public.users
		WHERE id=$1
	`

	var user models.User
	if err := ur.client.DbPool.QueryRow(ctx, query, id).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UserByUsername(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT *
		FROM public.users
		WHERE username=$1
	`

	var user models.User
	if err := ur.client.DbPool.QueryRow(ctx, query, id).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, id string, dto dto.UpdateUserDto) (*models.User, error) {
	query := `
		UPDATE public.users
		SET username=$1, password=$2
		WHERE id=$3
		RETURNING *
	`

	var user models.User
	if err := ur.client.DbPool.QueryRow(ctx, query, dt).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}