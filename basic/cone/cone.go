package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to a subject
	_, err = nc.Subscribe("pusher.new", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
