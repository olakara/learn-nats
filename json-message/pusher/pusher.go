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

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println("Error connecting to NATS:", err)
		return
	} else {
		fmt.Println("Connected to NATS")
	}
	defer nc.Close()

	stock := &stock{
		Symbol: "BTC",
		Price:  345.50,
	}

	payload, _ := json.Marshal(stock)
	err = nc.Publish("stock.update", payload)
	if err != nil {
		fmt.Println("Error publishing message to NATS:", err)
	} else {
		fmt.Printf("Published message to NATS at %s\n", time.Now().Format("15:04:05"))
	}

}
