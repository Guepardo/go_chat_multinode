package handlers

import (
	"log"

	"example.com/gochat/models"
)

type Middleware struct {
	chats map[string]map[*Websocket]bool

	register chan *Websocket

	unregister chan *Websocket

	exchangePublishChannel chan *models.Message

	exchangeConsumeChannel chan *models.Message
}

func NewMiddleware(exchangePublishChannel chan *models.Message, exchangeConsumeChannel chan *models.Message) *Middleware {
	return &Middleware{
		chats:                  make(map[string]map[*Websocket]bool),
		register:               make(chan *Websocket),
		unregister:             make(chan *Websocket),
		exchangePublishChannel: exchangePublishChannel,
		exchangeConsumeChannel: exchangeConsumeChannel,
	}
}

func (middleware *Middleware) handleSocketsLifeCycle() {
	for {
		select {
		case websocket := <-middleware.register:
			if _, ok := middleware.chats[websocket.ChatId]; !ok {
				middleware.chats[websocket.ChatId] = make(map[*Websocket]bool)
			}

			middleware.chats[websocket.ChatId][websocket] = true

			log.Default().Print("Novo login")

		case websocket := <-middleware.unregister:
			close(websocket.Send)
			delete(middleware.chats[websocket.ChatId], websocket)

			log.Default().Print("Novo logout")
		}
	}
}

func (middleware *Middleware) handleMessages() {
	for {
		select {
		case message := <-middleware.exchangeConsumeChannel:
			for websocket, _ := range middleware.chats[message.ChatId] {
				websocket.Send <- message.Text
			}
			log.Default().Print("Broadcast de mensagens")
		}
	}
}

func (middleware *Middleware) Start() {
	go middleware.handleSocketsLifeCycle()
	go middleware.handleMessages()
}
