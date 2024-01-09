package producer

import amqp "github.com/rabbitmq/amqp091-go"

type Channel interface {
	Close() error
	Confirm(noWait bool) error
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(c chan amqp.Return) chan amqp.Return
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}
