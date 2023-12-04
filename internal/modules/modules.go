package modules

import (
	"log/slog"
	"microauth/internal/modules/auth"
	"microauth/internal/storage/clients/postgresql"
)

type ModuleInitializer struct {
	Client *postgresql.PSQLClient
	Log    *slog.Logger
}

func New(client *postgresql.PSQLClient, log *slog.Logger) *ModuleInitializer {
	return &ModuleInitializer{
		Client: client,
		Log:    log,
	}
}

func (mi ModuleInitializer) CreateAllRoutes() {
	auth.New(mi.Client, mi.Log, AUTH_PATH)
}
