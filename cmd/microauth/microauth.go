package main

import (
	"fmt"
	"log/slog"
	"microauth/internal/domain/config"
	"microauth/internal/modules"
	"microauth/internal/storage/clients/postgresql"
	configManager "microauth/pkg/config"
	"microauth/pkg/logger/handlers/slogpretty"
	"net/http"
	"os"
	"sync"
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

	logger.Debug("[Postgres] Creating connection pool with database...")
	client, err := postgresql.New(cfg)
	if err != nil {
		panic("Cannot create connection pool with database. Error: " + err.Error())
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(clientGo *postgresql.PSQLClient, loggerGo *slog.Logger) {
		defer wg.Done()
		loggerGo.Debug("[HTTP and ModuleInitializer] Create connection and fetching all module routers...")
		RunServer(9090, clientGo, loggerGo)
	}(client, logger)

	wg.Wait()
}

func RunServer(port int, client *postgresql.PSQLClient, log *slog.Logger) {
	moduleInitializer := modules.New(client, log)
	moduleInitializer.CreateAllRoutes()

	//http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
	//	w.Write([]byte("Hello world!!!"))
	//})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
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
