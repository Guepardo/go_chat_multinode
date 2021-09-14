package main

import (
	"log"
	"net/http"

	"example.com/gochat/handlers"
)

func Start() {
	middleware := handlers.NewMiddleware()
	go middleware.Start()

	http.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
		handlers.StartWebsocket(middleware, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go WebChat!"))
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
