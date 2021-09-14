package handlers

import (
	"log"

	"example.com/gochat/models"
)

type Middleware struct {
	chats map[string]map[*Websocket]bool

	register chan *Websocket

	unregister chan *Websocket

	broadcast chan *models.Message
}

func NewMiddleware() *Middleware {
	return &Middleware{
		chats:      make(map[string]map[*Websocket]bool),
		register:   make(chan *Websocket),
		unregister: make(chan *Websocket),
		broadcast:  make(chan *models.Message),
	}
}

func (middleware *Middleware) Start() {
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

		case message := <-middleware.broadcast:
			for websocket, _ := range middleware.chats[message.ChatId] {
				websocket.Send <- message.Text
			}

			log.Default().Print("Broadcast de mensagens")
		}
	}
}
