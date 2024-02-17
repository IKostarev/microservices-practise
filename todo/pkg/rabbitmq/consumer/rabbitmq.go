package rabbitmq

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"todo/pkg/rabbitmq"
)

type RabbitMQ struct {
	connCfg *connectionConfig
	cfg     *rabbitmq.RabbitConsumerConfig
	logger  *zerolog.Logger

	ctxCancelFunc  context.CancelFunc
	batchConsumers map[string]*BatchConsumer
}

func New(cfg *rabbitmq.RabbitConsumerConfig, log *zerolog.Logger) (*RabbitMQ, error) {
	connCfg := &connectionConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
	}

	rmq := &RabbitMQ{
		cfg:            cfg,
		logger:         log,
		connCfg:        connCfg,
		batchConsumers: make(map[string]*BatchConsumer),
	}

	return rmq, nil
}

func (r *RabbitMQ) AddBatchConsumer(
	cfg *rabbitmq.RabbitConsumerConfig,
	queue string,
	exchange string,
	handler BatchHandler,
) error {
	consumer, err := newBatchConsumer(cfg, queue, exchange, handler, r.logger)
	if err != nil {
		return fmt.Errorf("create rabbitmq_consumer batch consumer: %w", err)
	}

	r.batchConsumers[queue] = consumer

	return nil
}

func (r *RabbitMQ) Run(ctx context.Context) {
	ctx, r.ctxCancelFunc = context.WithCancel(ctx)

	for _, batchConsumer := range r.batchConsumers {
		go batchConsumer.Consume(ctx)

		r.logger.Debug().Msgf(fmt.Sprintf("start consuming rabbitmq_consumer queue %s", batchConsumer.config.Queue))
	}
}

func (r *RabbitMQ) Stop() {
	r.ctxCancelFunc()
}
