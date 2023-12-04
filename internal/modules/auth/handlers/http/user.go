package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"io"
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

type ResponseForm[T ErrorForm | SignInForm | RegisterForm] struct {
	Code   int `json:"code"`
	Result T   `json:"result"`
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
	requestBody, _ := io.ReadAll(req.Body)
	var userDto *dto.UserDTO
	if err := json.Unmarshal(requestBody, userDto); err != nil {
		result, _ := json.Marshal(
			ResponseForm[ErrorForm]{
				Code:   int(http2.ErrCodeInternal),
				Result: ErrorForm{Message: fmt.Sprintf("Json-error %s", err.Error())},
			},
		)
		w.Write(result)
		return
	}
	token, err := auth.userService.SignIn(userDto)
	if err != nil {
		result, _ := json.Marshal(
			ResponseForm[ErrorForm]{
				Code:   int(http2.ErrCodeInternal),
				Result: ErrorForm{Message: fmt.Sprintf("SignIn-error %s", err.Error())},
			},
		)
		w.Write(result)
		return
	}

	result, _ := json.Marshal(
		ResponseForm[SignInForm]{
			Code:   200,
			Result: SignInForm{Token: token},
		},
	)
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
