package repositories

import (
	"context"
	"fmt"
	"microauth/internal/domain/models"
	"microauth/internal/modules/user/models/dto"
	"microauth/internal/storage/clients/postgresql"
	"time"
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
	if err := ur.client.DbPool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
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
	if err := ur.client.DbPool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	current, err := ur.UserById(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.Password == "" {
		dto.Password = current.PassHash
	}

	if dto.Username == "" {
		dto.Username = current.Username
	}
	var user models.User
	if err := ur.client.DbPool.QueryRow(ctx, query, dto.Username, dto.Password, id).Scan(&user.ID, &user.Username, &user.PassHash); err != nil {
		return nil, err
	}

	fmt.Println(user.Username, user.PassHash)
	return &user, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id string) (string, error) {
	query := `
		DELETE 
		FROM public.users
		WHERE id=$1
		RETURNING id
	`
	var deletedId string
	if err := ur.client.DbPool.QueryRow(ctx, query, &id).Scan(&deletedId); err != nil {
		return "", err
	}
	return deletedId, nil
}
