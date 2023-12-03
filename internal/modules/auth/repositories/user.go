package repositories

import (
	"context"
	"errors"
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
	`
	if err := ur.client.DbPool.QueryRow(ctx, query, ud.Username, ud.Password).Scan(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ur.log.Error(postgresql.GetError(err))
		}
		return pgErr
	}
	return nil
}

func (ur *UserRepository) User(ctx context.Context, ud *dto.UserDTO) (*models.User, error) {
	query := `
		SELECT id, username
		FROM public.users
		WHERE username=$1 AND password=$2
	`
	var userResult *models.User

	if err := ur.client.DbPool.QueryRow(ctx, query, ud.Username, ud.Password).Scan(&userResult.ID, &userResult.Username); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ur.log.Error(postgresql.GetError(err))
		}
		return nil, pgErr
	}
	return userResult, nil
}
