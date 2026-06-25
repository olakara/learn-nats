package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type stock struct {
	Symbol string  `json:"symbol"`
	Price  float32 `json:"price"`
}

func schduledJob(store *Store) {
	stock, ok := store.Pop()
	if !ok {
		fmt.Println("No stock items available to process.")
		return
	}
	fmt.Printf("Processing stock item: %+v\n", stock)
}

func main() {

	fmt.Println("Application started..")
	store := &Store{}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println("Error connecting to NATS:", err)
		return
	} else {
		fmt.Println("Connected to NATS")
	}
	defer nc.Close()

	sub, _ := nc.Subscribe("stock.update", func(msg *nats.Msg) {
		var data stock
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
		} else {
			fmt.Printf("Data received: %+v\n", data)
		}
		store.Append(data)
	})

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		schduledJob(store)
	}

	err = sub.AutoUnsubscribe(20)
	if err != nil {
		return
	}
	select {}

}
