package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"microauth/internal/domain/httpmanager"
	"microauth/internal/modules/auth/handlers"
	"microauth/internal/modules/auth/models/dto"
	"net/http"
)

type IAuthHandler interface {
	SignIn(w http.ResponseWriter, req *http.Request)
	//Register(w http.ResponseWriter, req *http.Request)
}

type AuthHandler struct {
	userService handlers.IUserService
}
type ErrorForm struct {
	Message string `json:"message"`
}
type SignInForm struct {
	Token string `json:"token"`
}
type RegisterForm struct {
	Status bool `json:"status"`
}

func (auth AuthHandler) SignIn(w http.ResponseWriter, req *http.Request) {
	var userDto dto.UserDTO

	if err := httpmanager.Request(req, &userDto); err != nil {
		result, err := httpmanager.Response[ErrorForm](
			http.StatusBadRequest,
			ErrorForm{Message: fmt.Sprintf("Json-error %s", err.Error())},
		)
		if err != nil {
			slog.Error("Global error")
			return
		}
		w.Write(result)
		return
	}

	token, err := auth.userService.SignIn(&userDto)
	if err != nil {
		result, err := httpmanager.Response[ErrorForm](
			http.StatusBadRequest,
			ErrorForm{Message: fmt.Sprintf("SignIn-error %s", err.Error())},
		)
		if err != nil {
			slog.Error("Global error")
			return
		}
		w.Write(result)
		return
	}

	result, _ := httpmanager.Response[SignInForm](http.StatusOK, SignInForm{Token: token})
	w.Write(result)
	return
}

func New(userService handlers.IUserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func AuthRouter(path string, authHandler IAuthHandler) {
	router := mux.NewRouter()

	fmt.Println(path)
	router.HandleFunc(fmt.Sprintf("/%s/signin", path), authHandler.SignIn).Methods("POST")
	//router.HandleFunc("/register", authHandler.RegisterHandle).Methods("POST")

	http.Handle("/", router)
}
