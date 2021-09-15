package exchanges

import (
	"example.com/gochat/models"
)

type StandaloneExchange struct {
	publish chan *models.Message
	consume chan *models.Message
}

func (s *StandaloneExchange) GetPublishChannel() chan *models.Message {
	return s.publish
}

func (s *StandaloneExchange) GetConsumeChannel() chan *models.Message {
	return s.consume
}

func (s *StandaloneExchange) Start() {
	for {
		select {
		case message := <-s.publish:
			s.consume <- message
		}
	}
}

func NewStandaloneExchange() *StandaloneExchange {
	return &StandaloneExchange{
		publish: make(chan *models.Message),
		consume: make(chan *models.Message),
	}
}
