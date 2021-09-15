package exchanges

import (
	"example.com/gochat/models"
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
	for {
		select {
		case message := <-s.publish:
			s.consume <- message
		}
	}
}

func NewRabbitMQExchange() *RabbitMQExchange {
	return &RabbitMQExchange{
		publish: make(chan *models.Message),
		consume: make(chan *models.Message),
	}
}
