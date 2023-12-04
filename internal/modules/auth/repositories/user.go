package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"microauth/internal/modules/auth/models/dto"
	"microauth/internal/modules/auth/models/models"
	"microauth/internal/storage/clients/postgresql"
)

type UserRepository struct {
	client *postgresql.PSQLClient
	log    *slog.Logger
}

func New(client *postgresql.PSQLClient) *UserRepository {
	return &UserRepository{client: client}
}

func (ur *UserRepository) AddUser(ctx context.Context, ud *dto.UserDTO) error {
	query := `
		INSERT INTO public.users
			(username, password)
		VALUES 
			($1, $2)
		RETURNING id
	`
	var id string
	if err := ur.client.DbPool.QueryRow(ctx, query, ud.Username, ud.Password).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		fmt.Println(err.Error() + " this is error")
		fmt.Println("Test3")
		if errors.As(err, &pgErr) {
			fmt.Println("Test4")
			ur.log.Error(postgresql.GetError(err))
		}
		fmt.Println("Test 1")
		return pgErr
	}

	fmt.Println("Test 2")
	return nil
}

func (ur *UserRepository) User(ctx context.Context, ud *dto.UserDTO) (*models.User, error) {
	query := `
		SELECT id, username
		FROM public.users
		WHERE username=$1 AND password=$2
	`
	userResult := models.User{ID: "", Username: "", PassHash: ""}

	if err := ur.client.DbPool.QueryRow(ctx, query, ud.Username, ud.Password).Scan(&userResult.ID, &userResult.Username, &userResult.PassHash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ur.log.Error(postgresql.GetError(err))
		}
		return nil, pgErr
	}
	return &userResult, nil
}

func (ur *UserRepository) UserWithoutPassword(ctx context.Context, ud *dto.UserDTO) (*models.User, error) {
	query := `
		SELECT id, username, password
		FROM public.users
		WHERE username=$1
	`
	userResult := models.User{ID: "", Username: "", PassHash: ""}

	if err := ur.client.DbPool.QueryRow(ctx, query, ud.Username).Scan(&userResult.ID, &userResult.Username, &userResult.PassHash); err != nil {
		// TODO: Адекватный вывод ошибки не работает
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ur.log.Error(postgresql.GetError(err))
		}
		return nil, pgErr
	}
	return &userResult, nil
}
