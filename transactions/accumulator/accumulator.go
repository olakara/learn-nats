package main

import (
	"accumulator/account"
	"accumulator/transactions"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println("Error connecting to NATS:", err)
		return
	} else {
		fmt.Println("Connected to NATS")
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		fmt.Println("Error getting JetStream context:", err)
		return
	}

	js.AddConsumer("accumulator", &nats.ConsumerConfig{
		Durable: "transactions-accumulator",
	})

	accountRepository := account.NewInMemoryAccountRepository()

	js.Subscribe("transactions.new", func(msg *nats.Msg) {
		var data transactions.Transaction
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
		}

		acc, err := accountRepository.GetByPersonId(data.PersonId)
		if err != nil {
			fmt.Println("No account for person ID:", data.PersonId)
			acc = account.NewAccount(uuid.New().String(), data.PersonId)
			accountRepository.CreateAccount(acc)
		}
		err = accountRepository.Accumulate(acc.Id, data.Amount)
		if err != nil {
			fmt.Println("Error accumulating amount:", err)
		}

		fmt.Printf("account %+v\n", acc)
	})

	select {}

}
