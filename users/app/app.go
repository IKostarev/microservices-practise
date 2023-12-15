package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"users/config"
	"users/internal/api"
	"users/internal/api/rest"
	"users/internal/repository"
	"users/internal/service"
	"users/pkg/logging"
	"users/pkg/postgresql"
)

type App struct {
	cfg         *config.Config
	logger      *zerolog.Logger
	router      *mux.Router
	userService api.UserService
}

func NewApp(
	cfg *config.Config,
) (*App, error) {
	// подключимся к базе данных
	databaseConn, err := postgresql.NewPgxConn(&cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// передадим подключение к базе данных констуктору репозитория
	userRepo := repository.NewUserRepository(databaseConn)

	// передадим реализацию репозитория конструктору сервиса
	userService := service.NewUserService(&cfg.Password, userRepo)

	logger := logging.NewLogger(cfg.Logging)

	return &App{
		cfg:         cfg,
		logger:      logger,
		userService: userService,
	}, nil
}

func (a *App) RunApp() {
	// инициализируем хэндлер
	userHandler := rest.NewUserHandler(a.logger, a.userService)

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

	// запустить вебсервер по адресу, передать в него роутер
	appAddr := fmt.Sprintf("%s:%s", a.cfg.App.AppHost, a.cfg.App.AppPort) // добавлен
	a.logger.Info().Msgf("running server at '%s'", appAddr)
	http.ListenAndServe(appAddr, a.router)
}
