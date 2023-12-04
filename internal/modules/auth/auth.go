package auth

import (
	"log/slog"
	"microauth/internal/modules/auth/handlers"
	authHttp "microauth/internal/modules/auth/handlers/http"
	"microauth/internal/modules/auth/repositories"
	"microauth/internal/modules/auth/services"
	"microauth/internal/storage/clients/postgresql"
)

type Auth struct {
	UserRepo    services.IUserRepo
	UserService handlers.IUserService
	Views       AuthViews
	log         *slog.Logger
}

type AuthViews struct {
	Http authHttp.IAuthHandler
}

func New(client *postgresql.PSQLClient, log *slog.Logger, httpRoute string) *Auth {
	userRepo := repositories.New(client)
	userService := services.New(userRepo, log)

	authHandler := authHttp.New(userService)
	if httpRoute != "" {
		authHttp.AuthRouter(httpRoute, authHandler)
	}
	return &Auth{
		UserRepo:    userRepo,
		UserService: userService,
		log:         log,
		Views:       AuthViews{Http: authHandler},
	}
}
