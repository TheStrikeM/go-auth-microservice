package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"microauth/internal/domain/httpmanager"
	"microauth/internal/domain/models"
	"microauth/internal/modules/user/handlers"
	"net/http"
)

type IUserHandler interface {
	UserById(w http.ResponseWriter, req *http.Request)
	UserByUsername(w http.ResponseWriter, req *http.Request)
	DeleteUser(w http.ResponseWriter, req *http.Request)
	UpdateUser(w http.ResponseWriter, req *http.Request)
}

type UserHandler struct {
	userService handlers.IUserService
}
type ErrorForm struct {
	Message string `json:"message"`
}
type UserByForm struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type DeleteUserForm struct {
	Status bool `json:"status"`
}
type UpdateUserForm struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserIdRequestForm struct {
	id string
}

func (uh *UserHandler) UserById(w http.ResponseWriter, req *http.Request) {
	var idForm UserIdRequestForm

	if err := httpmanager.Request(req, &idForm); err != nil {
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

	user, err := uh.userService.UserById(idForm.id)
	if err != nil {
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

	result, err := httpmanager.Response[models.User](
		http.StatusOK,
		*user,
	)
	w.Write(result)
	return
}

type UserUsernameRequestForm struct {
	username string
}

func (uh *UserHandler) UserByUsername(w http.ResponseWriter, req *http.Request) {
	var usernameForm UserUsernameRequestForm

	if err := httpmanager.Request(req, &usernameForm); err != nil {
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

	user, err := uh.userService.UserByUsername(usernameForm.username)
	if err != nil {
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

	result, err := httpmanager.Response[models.User](
		http.StatusOK,
		*user,
	)
	w.Write(result)
	return
}

func New(userService handlers.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func UserRouter(path string, userHandler UserHandler) {
	router := mux.NewRouter()

	fmt.Println(path)
	router.HandleFunc(fmt.Sprintf("/%s/get_by_id", path), userHandler.UserById).Methods("GET")
	router.HandleFunc(fmt.Sprintf("/%s/get_by_username", path), userHandler.UserByUsername).Methods("POST")

	http.Handle("/", router)
}
