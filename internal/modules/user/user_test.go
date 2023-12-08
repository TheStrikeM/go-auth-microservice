package user

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log/slog"
	"microauth/internal/domain/config"
	"microauth/internal/modules/user/models/dto"
	"microauth/internal/modules/user/repositories"
	"microauth/internal/modules/user/services"
	"microauth/internal/storage/clients/postgresql"
	configManager "microauth/pkg/config"
	"os"
	"testing"
)

func TestUserModuleUserById(t *testing.T) {
	cfg := configManager.MustLoad[config.Config]()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, _ := postgresql.New(cfg)

	userRepo := repositories.New(client)
	userService := services.New(userRepo, logger)

	value, err := userService.UserById("3d64d0f8-ebee-4542-9208-ce3249196792")
	require.NoError(t, err)
	fmt.Println(value)
}

func TestUserModuleUserByUsername(t *testing.T) {
	cfg := configManager.MustLoad[config.Config]()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, _ := postgresql.New(cfg)

	userRepo := repositories.New(client)
	userService := services.New(userRepo, logger)

	value, err := userService.UserByUsername("daria.so.05@mail.ru")
	require.NoError(t, err)
	fmt.Println(value)
}

func TestUserModuleDeleteUser(t *testing.T) {
	cfg := configManager.MustLoad[config.Config]()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, _ := postgresql.New(cfg)

	userRepo := repositories.New(client)
	userService := services.New(userRepo, logger)

	err := userService.DeleteUser("6e3408e3-de0a-4586-95b9-6603c9c91959")
	require.NoError(t, err)
}

func TestUserModuleUpdateUser(t *testing.T) {
	cfg := configManager.MustLoad[config.Config]()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, _ := postgresql.New(cfg)

	userRepo := repositories.New(client)
	userService := services.New(userRepo, logger)

	res, err := userService.UpdateUser("edb027da-716a-4a99-ba26-031b085b4423", &dto.UpdateUserDto{
		Username: "oleg",
	})
	require.NoError(t, err)

	fmt.Println(res)
}
