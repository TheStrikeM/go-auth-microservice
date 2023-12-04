package services

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"microauth/internal/modules/auth/models/dto"
	"microauth/internal/modules/auth/models/models"
	"microauth/pkg/hash"
	jwtManager "microauth/pkg/jwt"
	"time"
	"unicode"
)

var (
	ErrValidationAuth     = "Validation error"
	ErrUserExists         = "User already exists"
	ErrUserAdding         = "Error in AddUser"
	ErrInvalidCredentials = "Invalid credentials"
	ErrPasswordHash       = "Error in password hash"
)

type IUserRepo interface {
	AddUser(ctx context.Context, userDto *dto.UserDTO) error
	User(ctx context.Context, userDto *dto.UserDTO) (*models.User, error)
	UserWithoutPassword(ctx context.Context, userDto *dto.UserDTO) (*models.User, error)
}

type UserService struct {
	repo IUserRepo
	log  *slog.Logger
}

func New(userRepo IUserRepo, logger *slog.Logger) *UserService {
	return &UserService{
		repo: userRepo,
		log:  logger,
	}
}

func (us *UserService) SignIn(userDto *dto.UserDTO) (string, error) {
	if !isValidUsername(userDto.Username) || !isValidPassword(userDto.Password) {
		return "", fmt.Errorf("signin %s", ErrValidationAuth)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := us.repo.UserWithoutPassword(ctx, userDto)
	if err != nil {
		fmt.Println(user)
		return "", fmt.Errorf("signin %s", ErrInvalidCredentials)
	}
	if !hash.ComparePassword(user.PassHash, userDto.Password) {
		return "", fmt.Errorf("signin %s", ErrInvalidCredentials)
	}

	token, err := jwtManager.GenerateToken(
		jwt.MapClaims{
			"sub":      user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		},
	)

	if err != nil {
		return "", fmt.Errorf("signin %s", err.Error())
	}
	return token, nil
}

func (us *UserService) Register(userDto *dto.UserDTO) (err error) {
	// TODO: Customize validate errors
	if !isValidUsername(userDto.Username) || !isValidPassword(userDto.Password) {
		return fmt.Errorf("register %s", ErrValidationAuth)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := us.repo.UserWithoutPassword(ctx, userDto); err == nil {
		return fmt.Errorf("register %s", ErrUserExists)
	}

	userDto.Password, err = hash.HashPassword(userDto.Password)
	if err != nil {
		return fmt.Errorf("registe %s", ErrPasswordHash)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := us.repo.AddUser(ctx, userDto); err != nil {
		return fmt.Errorf("register %s is %s", ErrUserAdding, err.Error())
	}

	return nil
}

func isValidPassword(pass string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(pass) >= 10 {
		hasMinLen = true
	}
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
func isValidUsername(uname string) bool {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Var(uname, "required,email") == nil
}
