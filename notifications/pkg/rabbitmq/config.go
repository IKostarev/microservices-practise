package rabbitmq

type RabbitConfig struct {
	User            string `envconfig:"RABBITMQ_USER" default:"user"`
	Password        string `envconfig:"RABBITMQ_PASSWORD" default:"user"`
	Host            string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	Port            string `envconfig:"RABBITMQ_PORT" default:"5672"`
	MaxRetryAttempt int    `envconfig:"RABBITMQ_MAX_RETRY_ATTEMPT" default:"5"`
}
