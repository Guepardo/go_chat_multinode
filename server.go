package main

import (
	"log"
	"net/http"
	"os"

	"example.com/gochat/exchanges"
	"example.com/gochat/handlers"
)

func getExchange() exchanges.Exchange {
	exchangeType := os.Getenv("EXCHANGE_TYPE")

	if exchangeType == "rabbitmq" {
		return exchanges.NewRabbitMQExchange()
	}

	return exchanges.NewStandaloneExchange()
}

func Start() {
	exchange := getExchange()

	exchangePublishChannel := exchange.GetPublishChannel()
	exchangeConsumeChannel := exchange.GetConsumeChannel()

	middleware := handlers.NewMiddleware(exchangePublishChannel, exchangeConsumeChannel)

	go exchange.Start()
	go middleware.Start()

	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		handlers.StartWebsocket(middleware, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go WebChat!"))
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
