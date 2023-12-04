package auth

import (
	"log/slog"
	"microauth/internal/modules/auth/handlers"
	"microauth/internal/modules/auth/repositories"
	"microauth/internal/modules/auth/services"
	"microauth/internal/storage/clients/postgresql"
)

type Auth struct {
	UserRepo    services.IUserRepo
	UserService handlers.IUserService
	log         *slog.Logger
}

func New(client *postgresql.PSQLClient, log *slog.Logger) *Auth {
	userRepo := repositories.New(client)
	userService := services.New(userRepo, log)
	return &Auth{
		UserRepo:    userRepo,
		UserService: userService,
		log:         log,
	}
}
