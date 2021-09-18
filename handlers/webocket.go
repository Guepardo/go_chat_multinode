package handlers

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"example.com/gochat/models"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websocket struct {
	ChatId string

	middleware *Middleware

	conn *websocket.Conn

	Send chan []byte
}

func (w *Websocket) listenRead() {
	defer func() {
		w.middleware.unregister <- w
		w.conn.Close()
	}()

	w.conn.SetReadLimit(maxMessageSize)
	w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error { w.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, text, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		filteredText := bytes.TrimSpace(bytes.Replace(text, newline, space, -1))

		message := &models.Message{
			Text:   filteredText,
			ChatId: w.ChatId,
		}

		w.middleware.exchangePublishChannel <- message
	}
}

func (w *Websocket) listenWrite() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		w.conn.Close()
	}()

	for {
		select {
		case message, ok := <-w.Send:
			w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := w.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writer.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(w.Send)

			for i := 0; i < n; i++ {
				writer.Write(newline)
				writer.Write(<-w.Send)
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func StartWebsocket(middleware *Middleware, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	websocket := &Websocket{
		ChatId:     "1234",
		middleware: middleware,
		conn:       conn,
		Send:       make(chan []byte, 256),
	}

	websocket.middleware.register <- websocket

	go websocket.listenWrite()
	go websocket.listenRead()
}
