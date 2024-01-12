package producer

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"io"
	"sync"
	"time"
	"todo/pkg/rabbitmq"
)

type MsgData struct {
	Country string
	Data    []byte
}

type Connection struct {
	dsn              string
	exchangeName     string
	queueName        string
	reconnectTimeout time.Duration
	conn             *amqp.Connection
	serviceChannel   *amqp.Channel
	mu               sync.RWMutex
	isClosed         bool
}

type Producer struct {
	io.Closer
	Connection

	logger          *zerolog.Logger
	publishingChan  Channel
	maxRetryAttempt int
}

func New(
	cfg *rabbitmq.RabbitProducerConfig,
	exchangeName string,
	queueName string,
	log *zerolog.Logger,
) (*Producer, error) {
	rabbitDSN := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	ret := Producer{
		Connection: Connection{
			dsn:              rabbitDSN,
			exchangeName:     exchangeName,
			queueName:        queueName,
			reconnectTimeout: time.Second * 10,
		},

		logger:          log,
		maxRetryAttempt: cfg.MaxRetryAttempt,
	}
	err := ret.persistConnect()
	if err != nil {
		return nil, err
	}

	defer ret.serviceChannel.Close()

	return &ret, nil
}

func (p *Producer) connect() error {
	var err error

	p.conn, err = amqp.Dial(p.dsn)
	if err != nil {
		return err
	}

	if p.serviceChannel, err = p.MakeChannel(); err != nil {
		return err
	}

	err = p.serviceChannel.ExchangeDeclare(
		p.exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare amqp exchange %s: %w", p.exchangeName, err)
	}

	q, err := p.serviceChannel.QueueDeclare(
		p.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare amqp queue %s: %w", p.queueName, err)
	}

	err = p.serviceChannel.QueueBind(
		q.Name,
		p.queueName,
		p.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("bind amqp queue %s to %s: %w", q.Name, p.exchangeName, err)
	}

	if p.publishingChan, err = p.MakeChannel(); err != nil {
		return fmt.Errorf("make chan for amqp %s to %s: %w", q.Name, p.exchangeName, err)
	}

	p.logger.Info().Msg("successfully connected to rabbitMQ")
	return nil
}

func (p *Producer) persistConnect() error {
	if p.isClosed {
		return nil
	}
	if err := p.connect(); err != nil {
		return err
	}

	go func() {
		for {
			reason := <-p.conn.NotifyClose(make(chan *amqp.Error))
			if reason == nil {
				return
			}
			p.logger.Warn().Msgf("rabbitMQ connection closed: %s", reason.Reason)
			if p.isClosed {
				return
			}
			p.mu.Lock()
			for {
				if connErr := p.connect(); connErr != nil {
					p.logger.Warn().Msg("cannot establish rabbitMQ connection after timeout")
					time.Sleep(p.reconnectTimeout)
					continue
				}
				break
			}
			p.mu.Unlock()
		}
	}()

	return nil
}

func (p *Producer) Close() error {
	var ret error
	p.isClosed = true

	err := p.conn.Close()
	if err != nil {
		ret = multierror.Append(ret, err)
	}

	return ret
}

func (p *Producer) MakeChannel() (*amqp.Channel, error) {
	ch, err := p.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating amqp channel: %w", err)
	}

	return ch, nil
}

func (p *Producer) Publish(data []byte, requestID string) (err error) {
	for attempt := 1; attempt <= p.maxRetryAttempt; attempt++ {
		if p.publishingChan == nil {
			if p.publishingChan, err = p.MakeChannel(); err != nil {
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
			}
		}

		err = p.publishingChan.Publish(
			p.exchangeName,
			p.queueName,
			true,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         data,
				Headers:      amqp.Table{"requestId": requestID},
			})
		if err != nil {
			time.Sleep(time.Duration(attempt*attempt) * time.Second)
		} else {
			break
		}
	}

	return err
}
