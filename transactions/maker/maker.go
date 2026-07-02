package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maker/config"
	"maker/people"
	"maker/transaction"
	"os"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func main() {
	config := config.LoadConfig()
	fmt.Printf("Person Repo URL: %s\n", config.PersonURL)

	repo := people.NewInMemoryPersonRepository(config.PersonURL)
	fmt.Printf("Loaded %d people\n", len(repo.People()))

	if len(repo.People()) == 0 {
		fmt.Println("No people found. Please check the URL or the data source.")
		os.Exit(1)
	}

	p := repo.GetSinglePerson()
	fmt.Printf("Random Person: ID=%s, Name=%s\n", p.Id, p.Name)

	t := transaction.NewTransaction(uuid.New().String(), getRandomAmount(), p.Id)
	fmt.Printf("Created Transaction: ID=%s, Amount=%.2f, PersonID=%s\n", t.Id, t.Amount, t.PersonId)

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
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "transactions",
		Subjects: []string{"transactions.new"}, // can have more like "transactions.update" or "transactions.*" if needed
	})

	payload, _ := json.Marshal(t)
	if err := nc.Publish("transactions.new", payload); err != nil {
		fmt.Println("Error publishing transaction to NATS:", err)
	} else {
		fmt.Println("Transaction published to NATS")
	}

}

func getRandomAmount() float64 {
	return float64(0 + (uuid.New().ID() % 1000)) // Random amount between 0 and 999
}
