package rabbitmq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	exchangeName string
	conn         *amqp.Connection
	ch           *amqp.Channel
	q            amqp.Queue
}

func NewRabbitMQClient(host string, exchangeName string, queueName string) *RabbitMQClient {
	conn, err := amqp.Dial(host)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	rc := &RabbitMQClient{
		conn:         conn,
		ch:           ch,
		exchangeName: exchangeName,
	}
	rc.QueueDeclare(queueName)
	return rc
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (r *RabbitMQClient) QueueDeclare(name string) {
	q, err := r.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	r.q = q
}

type RabbitMQProducer struct {
	RabbitMQClient
}

func NewRabbitMQProducer(client *RabbitMQClient) *RabbitMQProducer {
	p := &RabbitMQProducer{
		RabbitMQClient: *client,
	}
	return p
}

func (r *RabbitMQProducer) Publish(file string) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	var jsonArr []interface{}
	err = json.Unmarshal(data, &jsonArr)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	for _, obj := range jsonArr {
		data, err = json.Marshal(obj)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		select {
			case <-r.ch.NotifyClose()
		}
		
		fmt.Println(r.conn.IsClosed())
		// defer r.conn.Close()
		err := r.ch.Publish(
			r.exchangeName, // exchange
			r.q.Name,       // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(data),
			})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
