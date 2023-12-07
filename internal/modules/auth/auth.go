package auth

import (
	"log/slog"
	"microauth/internal/modules/auth/handlers"
	authHttp "microauth/internal/modules/auth/handlers/http"
	"microauth/internal/modules/auth/repositories"
	"microauth/internal/modules/auth/services"
	"microauth/internal/storage/clients/postgresql"
)

type AuthModule struct {
	Repo    services.IUserRepo
	Service handlers.IUserService
	Views   Views
	log     *slog.Logger
}

type Views struct {
	Http authHttp.IAuthHandler
}

func New(client *postgresql.PSQLClient, log *slog.Logger, httpRoute string) *AuthModule {
	userRepo := repositories.New(client)
	userService := services.New(userRepo, log)

	authHandler := authHttp.New(userService)
	if httpRoute != "" {
		authHttp.AuthRouter(httpRoute, authHandler)
	}
	return &AuthModule{
		Repo:    userRepo,
		Service: userService,
		log:     log,
		Views:   Views{Http: authHandler},
	}
}
