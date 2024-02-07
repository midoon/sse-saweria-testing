package main

import (
	"fmt"
	"net/http"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "failed to build connection", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		rmqConn, err := rmqConsumerInit()
		if err != nil {
			panic(err)
		}

		defer rmqConn.Close()

		rmqChannel, err := rmqConn.Channel()
		if err != nil {
			panic(err)
		}

		saweriaConsumer, err := rmqChannel.ConsumeWithContext(r.Context(), "saweria", "saweria-consumer", true, false, false, false, nil)
		if err != nil {
			panic(err)
		}

		for message := range saweriaConsumer {
			event := fmt.Sprintf("event: %s\n"+"data: %s \n\n", "sawer-kuy", string(message.Body))
			_, _ = fmt.Fprint(w, event)
			flusher.Flush()
		}
	})

	server := http.Server{
		Addr:    "localhost:5001",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func rmqConsumerInit() (*amqp091.Connection, error) {
	rwmqConnection, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	return rwmqConnection, err
}
