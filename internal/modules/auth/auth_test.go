package auth

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log/slog"
	"microauth/internal/domain/config"
	"microauth/internal/modules/auth/models/dto"
	"microauth/internal/storage/clients/postgresql"
	configManager "microauth/pkg/config"
	"os"
	"testing"
)

func TestSignIn(t *testing.T) {
	userDto := dto.UserDTO{
		Username: "thestrikem@yander.ru",
		Password: "valentin228Tubik@s",
	}
	cfg := configManager.MustLoad[config.Config]()

	client, err := postgresql.New(cfg)
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	require.NoError(t, err)

	authModule := New(client, log, "auth")
	token, err := authModule.UserService.SignIn(&userDto)
	require.NoError(t, err)
	fmt.Println(token)
}

func TestRegister(t *testing.T) {
	userDto := dto.UserDTO{
		Username: "thestrikem@yander.ru",
		Password: "valentin228Tubik@s",
	}
	cfg := configManager.MustLoad[config.Config]()

	client, err := postgresql.New(cfg)
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	require.NoError(t, err)

	authModule := New(client, log, "auth")
	err = authModule.UserService.Register(&userDto)
	require.NoError(t, err)
}
