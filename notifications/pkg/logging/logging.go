package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type LoggerConfig struct {
	LogIndex  string `envconfig:"LOG_INDEX" required:"true" default:"my_app"`
	IsDebug   bool   `envconfig:"LOG_IS_DEBUG" required:"true" default:"false"`
	LogToFile bool   `envconfig:"LOG_TO_FILE" required:"true" default:"true"`
}

func NewLogger(cfg LoggerConfig) *zerolog.Logger {
	// Создадим новый логгер с полем "service" и в него будем писать LogIndex,
	// чтобы понимать, какому из сервисов принадлежит сообщение.
	logger := log.With().Str("service", cfg.LogIndex).Logger()

	// Установим уровень логирования на основе IsDebug.
	if cfg.IsDebug {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}

	// Опеределим. куда будем писать наши логи - os.STDOUT или в файл,
	// в зависимости от LogToFile.
	if cfg.LogToFile {
		// В файл
		file, err := os.Create("application.log")
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create log file")
		}
		logger = logger.Output(file)
	} else {
		// В os.STDOUT
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return &logger
}
