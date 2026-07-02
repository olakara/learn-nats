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

	t := transaction.NewTransaction(uuid.New().String(), 100.0, p.Id)
	fmt.Printf("Created Transaction: ID=%s, Amount=%.2f, PersonID=%s\n", t.Id, t.Amount, t.PersonId)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	payload, _ := json.Marshal(t)
	if err := nc.Publish("transactions.new", payload); err != nil {
		log.Fatal(err)
	}

}
