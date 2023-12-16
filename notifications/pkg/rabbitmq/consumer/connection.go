package rabbitmq

import (
	"fmt"
	"github.com/rs/zerolog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	logger     *zerolog.Logger
	cfg        *connectionConfig
	connection *amqp.Connection
	channel    *amqp.Channel
}

type connectionConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewConnection(config *connectionConfig, log *zerolog.Logger) (*Connection, error) {
	c := &Connection{
		cfg:    config,
		logger: log,
	}

	err := c.connect()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) connect() error {
	connection, err := amqp.Dial(c.url())
	if err != nil {
		return fmt.Errorf("amqp dial: %w", err)
	}
	c.connection = connection

	return c.connectChannel()
}

func (c *Connection) reconnect() {
	reconnectFunc := c.connect
	if c.connection != nil && !c.connection.IsClosed() {
		reconnectFunc = c.connectChannel
	}

	attempt := 1
	for {
		err := reconnectFunc()
		if err == nil {
			c.logger.Info().Msg("reconnected to rabbitmq_consumer")

			return
		}

		c.logger.Warn().Msg("failed to establish rabbitmq_consumer connection, retrying")

		time.Sleep(time.Duration(attempt*attempt) * time.Second)
		attempt++
	}
}

func (c *Connection) connectChannel() error {
	channel, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("create amqp channel %w", err)
	}
	c.channel = channel

	return nil
}

func (c *Connection) url() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.cfg.User, c.cfg.Password, c.cfg.Host, c.cfg.Port)
}

func (c *Connection) Close() error {
	return c.connection.Close()
}
