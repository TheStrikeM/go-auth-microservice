package main

import (
	"fmt"
	"log/slog"
	"microauth/internal/domain/config"
	"microauth/internal/storage/clients/postgresql"
	configManager "microauth/pkg/config"
	"microauth/pkg/logger/handlers/slogpretty"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	fmt.Println("Starting loading, opening config...")
	cfg := configManager.MustLoad[config.Config]()

	fmt.Println("Start loading logger...")
	logger := setupLogger(cfg.Env)
	logger.Debug("Yey! Logger, Config enabled!")

	logger.Debug("Creating connection pool with database...")
	_, err := postgresql.New(cfg)
	if err != nil {
		panic("Cannot create connection pool with database. Error: " + err.Error())
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
