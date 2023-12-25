package rabbitmq

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"users/pkg/rabbitmq"
)

type BatchHandler interface {
	Handle(context.Context, []amqp.Delivery) (uint64, []Error, error)
}

type batchConsumerConfig struct {
	Queue         string
	Exchange      string
	BatchSize     int
	BatchWaitTime time.Duration
}

type BatchConsumer struct {
	logger  *zerolog.Logger
	config  *batchConsumerConfig
	conn    *Connection
	handler BatchHandler
}

type Error struct {
	Msg amqp.Delivery
	Err error
}

func newBatchConsumer(
	cfg *rabbitmq.RabbitConsumerConfig,
	queue string,
	exchange string,
	handler BatchHandler,
	logger *zerolog.Logger,
) (*BatchConsumer, error) {
	connCfg := &connectionConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
	}

	conn, err := NewConnection(connCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq_consumer new connection fail %w", err)
	}

	consumerCfg := &batchConsumerConfig{
		Queue:         queue,
		Exchange:      exchange,
		BatchSize:     cfg.BatchSize,
		BatchWaitTime: cfg.BatchWaitTime,
	}

	consumer := &BatchConsumer{
		handler: handler,
		conn:    conn,
		config:  consumerCfg,
		logger:  logger,
	}

	if err := consumer.connect(); err != nil {
		return nil, fmt.Errorf("rabbitmq_consumer consumer connect fail %w", err)
	}

	return consumer, nil
}

func (bc *BatchConsumer) connect() error {
	cfg := bc.config
	ch := bc.conn.channel

	err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare amqp exchange %s: %w", cfg.Exchange, err)
	}

	queue, err := ch.QueueDeclare(cfg.Queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare amqp queue %s: %w", cfg.Queue, err)
	}

	err = ch.QueueBind(queue.Name, cfg.Queue, cfg.Exchange, false, nil)
	if err != nil {
		return fmt.Errorf("bind amqp queue %s to exchange %s: %w", cfg.Queue, cfg.Exchange, err)
	}

	err = ch.Qos(cfg.BatchSize, 0, true)
	if err != nil {
		return fmt.Errorf("amqp qos: %w", err)
	}

	bc.logger.Info().Msgf("connection to rabbitmq for %s queue established", cfg.Queue)

	return nil
}

func (bc *BatchConsumer) Consume(ctx context.Context) {
	messages := bc.delivery()
	buffer := bc.newBatchBuffer()
	ticker := time.NewTicker(bc.config.BatchWaitTime)

	for {
		select {
		case <-ctx.Done():
			bc.logger.Debug().Msgf("rabbitmq consumer for queue %s is stopped by context signal", bc.QueueName())
			bc.processBatch(ctx, buffer)

			bc.conn.Close()
			return
		case <-ticker.C:
			bc.logger.Debug().Msgf("processing messages in queue %s by ticker signal", bc.QueueName())

			bc.processBatch(ctx, buffer)
			buffer = bc.newBatchBuffer()
		case msg, ok := <-messages:
			if !ok {
				buffer = bc.newBatchBuffer()

				bc.conn.reconnect()
				messages = bc.delivery()

				continue
			}

			buffer = append(buffer, msg)
			if len(buffer) < bc.config.BatchSize {
				continue
			}

			bc.processBatch(ctx, buffer)
			buffer = bc.newBatchBuffer()
		}
	}
}

func (bc *BatchConsumer) delivery() <-chan amqp.Delivery {
	delivery, err := bc.conn.channel.Consume(
		bc.QueueName(), "", false, false, false, false, nil,
	)
	if err != nil {
		bc.logger.Error().Msgf(
			"consume failed", zap.Error(err),
			zap.String("queue", bc.QueueName()),
		)

		return nil
	}

	return delivery
}

func (bc *BatchConsumer) processBatch(ctx context.Context, messages []amqp.Delivery) {
	if len(messages) == 0 {
		return
	}

	ch := bc.conn.channel

	lastSuccessTag, handlerErrList, fatalErr := bc.handler.Handle(ctx, messages)
	if fatalErr != nil {
		bc.logger.Error().Msgf("batch handle failed", zap.Error(fatalErr))

		lastMessage := messages[len(messages)-1]
		if err := ch.Nack(lastMessage.DeliveryTag, true, true); err != nil {
			bc.logger.Error().Msgf("nack failed", zap.Error(err))
		}

		return
	}

	for _, consumeErr := range handlerErrList {
		bc.logger.Error().Msgf(
			fmt.Sprintf("handle failed in queue %s", bc.QueueName()),
			zap.Error(consumeErr.Err),
			zap.ByteString("rabbit message", consumeErr.Msg.Body),
		)
		if err := ch.Reject(consumeErr.Msg.DeliveryTag, false); err != nil {
			bc.logger.Error().Msgf("reject failed", zap.Error(err))
		}
	}

	if lastSuccessTag > 0 {
		if err := ch.Ack(lastSuccessTag, true); err != nil {
			bc.logger.Error().Msgf("ack failed", zap.Error(err))
		}

		bc.logger.Debug().Msgf("processed batch from queue %s, last success tag %d", bc.QueueName(), lastSuccessTag)
	}
}

func (bc *BatchConsumer) QueueName() string {
	return bc.config.Queue
}

func (bc *BatchConsumer) newBatchBuffer() []amqp.Delivery {
	return make([]amqp.Delivery, 0, bc.config.BatchSize)
}
