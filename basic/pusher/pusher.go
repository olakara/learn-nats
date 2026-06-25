package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func createMessage() string {
	now := time.Now()
	return fmt.Sprintf("Message at %s\n", now.Format("2006-01-02 15:04:05"))
}

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println("Error connecting to NATS:", err)
		return
	} else {
		fmt.Println("Connected to NATS")
	}

	defer nc.Close()

	for {
		message := createMessage()
		fmt.Println(message)
		err = nc.Publish("pusher.new", []byte(message))
		if err != nil {
			fmt.Println("Error publishing message:", err)
		}
		time.Sleep(1 * time.Second)
	}
}
