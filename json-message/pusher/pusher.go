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

func getNewStockPrice() float32 {
	// Generate a random stock price between 300 and 400
	return 300 + float32(time.Now().UnixNano()%100)/10
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

	//create a random stock price and publish it to the NATS server every 5 seconds
	for {
		stock := stock{
			Symbol: "AAPL",
			Price:  getNewStockPrice(),
		}

		payload, _ := json.Marshal(stock)
		err = nc.Publish("stock.update", payload)
		if err != nil {
			fmt.Println("Error publishing message to NATS:", err)
		} else {
			fmt.Printf("Published message to NATS at %s\n", time.Now().Format("15:04:05"))
		}
		time.Sleep(2 * time.Second)
	}

}
