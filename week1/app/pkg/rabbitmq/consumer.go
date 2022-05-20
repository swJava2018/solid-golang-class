package rabbitmq

import (
	"encoding/json"
	"errors"
	"event-data-pipeline/pkg/logger"
	"time"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Create() error
	Connect() error
	CreateChannel() error
	Delete() error
	ReConnect() error
	FetchRecords() (map[string]interface{}, error)
	ExchangeDeclare() error
	QueueBind() error
	InitDeliveryChannel() error
}

type RabbitMQConsumer struct {
	host           string
	exchangeName   string
	exchangeType   string
	queueName      string
	routingKey     string
	conn           *amqp.Connection
	ConnMaxRetries int
	ch             *amqp.Channel
	q              amqp.Queue
	message        <-chan amqp.Delivery
	err            chan error
	shutdown       chan bool
}

func NewRabbitMQConsumer(host string, exchangeName string, exchangeType string, queueName string, routingKey string) *RabbitMQConsumer {

	c := &RabbitMQConsumer{
		host:         host,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queueName,
		routingKey:   routingKey,
		err:          make(chan error),
		shutdown:     make(chan bool, 1),
	}
	return c
}
func (c *RabbitMQConsumer) Create() error {

	// best practice is to reuse connections and channels
	err := c.Connect()
	if err != nil {
		return err
	}
	err = c.CreateChannel()
	if err != nil {
		return err
	}
	logger.Debugf("Check in consumer creation: %v", c.ch)

	return nil
}

func (c *RabbitMQConsumer) Connect() error {
	c.ConnMaxRetries = -1
	retry := 0
	for {
		conn, err := amqp.Dial(c.host)
		if conn != nil {
			logger.Infof("rabbitmq connection made")
			c.conn = conn
			break
		}
		retry++
		if c.ConnMaxRetries >= 0 && retry > c.ConnMaxRetries {
			return err
		}
		time.Sleep(1 * time.Second)
		logger.Infof("retry to connect to rabbitmq")
	}
	//Listen to NotifyClose
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error))
		c.err <- errors.New("rabbitmq connection closed")
	}()
	return nil
}

func (c *RabbitMQConsumer) CreateChannel() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	c.ch = ch
	return nil
}

func (c *RabbitMQConsumer) Delete() error {
	logger.Debugf("deleting rabbit mq consumer connection: %s and channel: %s", c.conn, c.ch)
	if c.ch != nil {
		err := c.ch.Close()
		if err != nil {
			return err
		}
	}
	if c.conn != nil && !c.conn.IsClosed() {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *RabbitMQConsumer) ReConnect() error {
	logger.Debugf("Reconnecting")
	err := c.Connect()
	if err != nil {
		return err
	}
	err = c.CreateChannel()
	if err != nil {
		return err
	}

	err = c.InitDeliveryChannel()
	if err != nil {
		return err
	}
	return nil

}

func (c *RabbitMQConsumer) FetchRecords() (map[string]interface{}, error) {
	var body map[string]interface{}

	// Check RabbitMQ Connection Error
	if c.conn == nil {
		err := c.ReConnect()
		if err != nil {
			return nil, err
		}
	}

	logger.Debugf("waiting letter")

	var letter amqp.Delivery
	select {
	case message := <-c.message:
		letter = message
	case <-c.shutdown:
		return nil, nil
	}
	logger.Debugf("Recevied Letter :%s", string(letter.Body))
	err := json.Unmarshal(letter.Body, &body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *RabbitMQConsumer) ExchangeDeclare() error {

	err := c.ch.ExchangeDeclare(
		c.exchangeName, // name
		c.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	return nil
}
func (c *RabbitMQConsumer) QueueDeclare() error {
	q, err := c.ch.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	c.q = q
	if err != nil {
		return err
	}
	return nil
}

func (c *RabbitMQConsumer) QueueBind() error {
	err := c.ch.QueueBind(
		c.queueName,    // queue name
		c.routingKey,   // routing key
		c.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *RabbitMQConsumer) InitDeliveryChannel() error {
	msg, err := c.ch.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}
	c.message = msg
	return nil
}
