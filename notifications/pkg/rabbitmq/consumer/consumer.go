package consumer

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"io"
	"notifications/pkg/rabbitmq"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MsgData struct {
	Country string
	Data    []byte
}

type handlerConfig struct {
	queue    string
	exchange string
	handler  Handler
}

type Connection struct {
	dsn              string
	reconnectTimeout time.Duration
	conn             *amqp.Connection
	serviceChannel   *amqp.Channel
	mu               sync.RWMutex
	isClosed         bool
}

type Consumer struct {
	io.Closer
	Connection

	logger           *zerolog.Logger
	consumingChannel Channel
	maxRetryAttempt  int
	handlers         map[string]handlerConfig
}

func New(
	cfg *rabbitmq.RabbitConfig,
	log *zerolog.Logger,
) (*Consumer, error) {
	rabbitDSN := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	ret := Consumer{
		Connection: Connection{
			dsn:              rabbitDSN,
			reconnectTimeout: time.Second * 10,
		},
		logger:          log,
		maxRetryAttempt: cfg.MaxRetryAttempt,
		handlers:        make(map[string]handlerConfig),
	}
	err := ret.persistConnect()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (c *Consumer) connect() error {
	var err error

	c.conn, err = amqp.Dial(c.dsn)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) persistConnect() error {
	if c.isClosed {
		return nil
	}
	if err := c.connect(); err != nil {
		return err
	}

	go func() {
		for {
			reason := <-c.conn.NotifyClose(make(chan *amqp.Error))
			if reason == nil {
				return
			}
			c.logger.Warn().Msgf("rabbitMQ connection closed: %s", reason.Reason)
			if c.isClosed {
				return
			}
			c.mu.Lock()
			for {
				if connErr := c.connect(); connErr != nil {
					c.logger.Warn().Msg("cannot establish rabbitMQ connection after timeout")
					time.Sleep(c.reconnectTimeout)
					continue
				}
				break
			}
			c.mu.Unlock()
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	var ret error
	c.isClosed = true

	err := c.conn.Close()
	if err != nil {
		ret = multierror.Append(ret, err)
	}

	return ret
}

func (c *Consumer) MakeChannel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating amqp channel: %w", err)
	}

	return ch, nil
}

func (c *Consumer) SetHandler(queueName, exchangeName string, handler Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.handlers == nil {
		c.handlers = make(map[string]handlerConfig)
	}

	c.handlers[queueName] = handlerConfig{
		queue:    queueName,
		exchange: exchangeName,
		handler:  handler,
	}
}

func (c *Consumer) Consume(queueName, exchangeName, exchangeType string, queueDurable, exchangeDurable bool) error {
	go func() {
		c.mu.Lock()
		ch, err := c.MakeChannel()
		c.mu.Unlock()
		if err != nil {
			c.logger.Error().Msgf("Error making channel: %s", err)
			return
		}

		err = ch.ExchangeDeclare(
			exchangeName,
			exchangeType,
			exchangeDurable,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			c.logger.Error().Msgf("Error declaring exchange: %s", err)
			return
		}

		_, err = ch.QueueDeclare(
			queueName,
			queueDurable,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			c.logger.Error().Msgf("Error declaring queue: %s", err)
			return
		}

		err = ch.QueueBind(
			queueName,
			"",
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			c.logger.Error().Msgf("Error binding queue: %s", err)
			return
		}

		msgs, err := ch.Consume(
			queueName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			c.logger.Error().Msgf("Error consuming messages: %s", err)
			return
		}

		for msg := range msgs {
			if handler, ok := c.handlers[queueName]; ok {
				handler.handler.Handle(msg)
			} else {
				c.logger.Error().Msgf("No handler found for queue %s", queueName)
			}
		}
	}()
	return nil
}

func (c *Consumer) Run() {
	for _, cfg := range c.handlers {
		err := c.Consume(cfg.queue, cfg.exchange, "direct", true, true)
		if err != nil {
			c.logger.Error().Msgf("Error setting up consumer for queue %s: %s", cfg.queue, err)
		}
	}
}
