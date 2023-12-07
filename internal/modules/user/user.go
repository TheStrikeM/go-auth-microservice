package user

import (
	"log/slog"
	authHttp "microauth/internal/modules/auth/handlers/http"
	"microauth/internal/modules/auth/repositories"
	"microauth/internal/modules/auth/services"
	"microauth/internal/storage/clients/postgresql"
)

type UserModule struct {
	Repo    string
	Service string
	Views   Views
	log     *slog.Logger
}

type Views struct {
	Http authHttp.IAuthHandler
}

func New(client *postgresql.PSQLClient, log *slog.Logger, httpRoute string) *UserModule {
	userRepo := repositories.New(client)
	userService := services.New(userRepo, log)

	authHandler := authHttp.New(userService)
	if httpRoute != "" {
		authHttp.AuthRouter(httpRoute, authHandler)
	}
	return &UserModule{
		Repo:    "userRepo",
		Service: "userService",
		log:     log,
		Views:   Views{Http: authHandler},
	}
}
