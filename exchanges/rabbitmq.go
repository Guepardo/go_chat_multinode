package exchanges

import (
	"log"
	"os"

	"example.com/gochat/models"
	"github.com/streadway/amqp"
)

const (
	exchangeName = "msg.exchange"
	routingKey   = "msg"
)

type RabbitMQExchange struct {
	publish chan *models.Message
	consume chan *models.Message
}

func (s *RabbitMQExchange) GetPublishChannel() chan *models.Message {
	return s.publish
}

func (s *RabbitMQExchange) GetConsumeChannel() chan *models.Message {
	return s.consume
}

func (s *RabbitMQExchange) Start() {
	log.Default().Print("Exchange Type: RabbitMQ")

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	connection, err := amqp.Dial(amqpServerURL)

	if err != nil {
		panic(err)
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	queue, err := channel.QueueDeclare(
		"",    // queue name
		false, // durable
		true,  // auto delete
		true,  // exclusive
		false, // no wait
		nil,   // arguments
	)

	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil)

	if err != nil {
		panic(err)
	}

	consumeChannel, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		panic(err)
	}

	for {
		select {
		case message := <-s.publish:
			channel.Publish(
				exchangeName,
				routingKey,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(message.Text),
				},
			)
		case message := <-consumeChannel:
			s.consume <- &models.Message{
				Text:   message.Body,
				ChatId: "1234",
			}
		}
	}
}

func NewRabbitMQExchange() *RabbitMQExchange {
	return &RabbitMQExchange{
		publish: make(chan *models.Message),
		consume: make(chan *models.Message),
	}
}
