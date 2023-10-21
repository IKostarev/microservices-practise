package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"users/internal/handlers"
	"users/internal/service"
	"users/pkg/jwtutil"
)

type App struct {
	router      *mux.Router
	userService handlers.UserService
}

func NewApp() (*App, error) {
	userService := service.NewUserService(nil, jwtutil.JWTUtil{})

	return &App{
		userService: userService,
	}, nil
}

func (a *App) RunApp() {
	// инициализируем хэндлер
	userHandler := handlers.NewUserHandler(a.userService)

	// инициализация роутера и сохранение его в соотвтетсвующее поле приложения
	a.router = mux.NewRouter()

	// зарегистрировать нового пользователя
	a.router.HandleFunc("/users/register", userHandler.RegisterUser).Methods(http.MethodPost)
	// получить пользователя по айди
	a.router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetGetUserById).Methods(http.MethodGet)
	// обновить пользователя
	a.router.HandleFunc("/users/update", userHandler.UpdateUser).Methods(http.MethodPut)
	// обновить пароль
	a.router.HandleFunc("/users/update-password", userHandler.UpdatePassword).Methods(http.MethodPut)
	// удалить пользователя
	a.router.HandleFunc("/users/delete/{id:[0-9]+}", userHandler.DeleteUser).Methods(http.MethodDelete)

	// выполнить вход в систему
	a.router.HandleFunc("/users/login", userHandler.UserLogin).Methods(http.MethodPost)
	// обновить токены пользователя
	a.router.HandleFunc("/users/refresh", userHandler.Refresh).Methods(http.MethodPost)
	// верифицировать токены пользователя
	a.router.HandleFunc("/users/verify", userHandler.Verify).Methods(http.MethodPost)

	// запустить вебсервер по адресу, передать в него роутер
	http.ListenAndServe(":8080", a.router)
}
