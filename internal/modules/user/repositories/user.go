package repositories

import (
	"context"
	"microauth/internal/domain/models"
	"microauth/internal/modules/user/models/dto"
	"microauth/internal/storage/clients/postgresql"
)

type UserRepository struct {
	client *postgresql.PSQLClient
}

func New(client *postgresql.PSQLClient) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
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

func (ur *UserRepository) UpdateUser(ctx context.Context, id string, dto *dto.UpdateUserDto) (*models.User, error) {
	query := `
		UPDATE public.users
		SET username=$1, password=$2
		WHERE id=$3
		RETURNING *
	`

	var user models.User
	if err := ur.client.DbPool.QueryRow(ctx, query, dto.Username, dto.Password, id).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		DELETE 
		FROM public.users
		WHERE  id=$1
	`

	if err := ur.client.DbPool.QueryRow(ctx, query).Scan(); err != nil {
		return err
	}

	return nil
}
