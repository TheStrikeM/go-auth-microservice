package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"microauth/internal/domain/config"
	"microauth/internal/modules"
	"microauth/internal/storage/clients/postgresql"
	configManager "microauth/pkg/config"
	"microauth/pkg/logger/handlers/slogpretty"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	wg.Add(1)
	go func(ctx context.Context, clientGo *postgresql.PSQLClient, loggerGo *slog.Logger) {
		defer wg.Done()
		loggerGo.Debug("[HTTP] Create connection and fetching all module routers...")
		MustRunServer(ctx, 9090, clientGo, loggerGo)
	}(ctx, client, logger)
	wg.Wait()

}

func MustRunServer(ctx context.Context, port int, client *postgresql.PSQLClient, log *slog.Logger) {
	moduleInitializer := modules.New(client, log)
	moduleInitializer.CreateAllRoutes()

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Debug("[HTTP] Enabling server, please wait...")
		return httpServer.ListenAndServe()
	})

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	g.Go(func() error {
		log.Debug("[HTTP] Start handling signals")
		<-gCtx.Done()
		log.Debug("[HTTP] Disabling server, please wait...")
		return httpServer.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("[HTTP] Server success closed, good luck")
		}
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
