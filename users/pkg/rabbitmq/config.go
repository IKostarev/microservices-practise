package rabbitmq

import "time"

type RabbitBaseConfig struct {
	User     string `envconfig:"RABBITMQ_USER" default:"user"`
	Password string `envconfig:"RABBITMQ_PASSWORD" default:"user"`
	Host     string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	Port     string `envconfig:"RABBITMQ_PORT" default:"5672"`
}

type RabbitConsumerConfig struct {
	RabbitBaseConfig
	BatchSize     int           `envconfig:"RABBITMQ_BATCH_SIZE" default:"100"`
	BatchWaitTime time.Duration `envconfig:"RABBITMQ_BATCH_WAIT_TIME" default:"10s"`
}

type RabbitProducerConfig struct {
	RabbitBaseConfig
	MaxRetryAttempt int `envconfig:"RABBITMQ_MAX_RETRY_ATTEMPT" default:"5"`
}
