package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"microauth/internal/domain/config"
	"time"
)

type PSQLClient struct {
	DbPool *pgxpool.Pool
}

func New(c *config.Config) (client *PSQLClient, err error) {
	var dbPool *pgxpool.Pool

	if err := DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := createUrl(c.Postgres)
		dbPool, err = pgxpool.New(ctx, url)
		if err != nil {
			return err
		}

		return nil
	}, 4, 2*time.Second); err != nil {
		return nil, err
	}

	return &PSQLClient{
		DbPool: dbPool,
	}, nil
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}

	return
}

func createUrl(cfg config.PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Db,
	)
}

func GetError(err error) string {
	pgErr := err.(*pgconn.PgError)
	return fmt.Sprintf(
		"SQL ERROR - %s, Detail - %s, Where - %s, Code - %s, SQLState - %s",
		pgErr.Message,
		pgErr.Detail,
		pgErr.Where,
		pgErr.Code,
		pgErr.SQLState(),
	)
}
