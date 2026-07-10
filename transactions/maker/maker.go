package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maker/config"
	"maker/people"
	"maker/tracing"
	"maker/transaction"
	"os"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx := context.Background()

	shutdown, err := tracing.Init(ctx, "maker")
	if err != nil {
		log.Fatal("Error initializing tracing:", err)
	}
	// Flush buffered spans before exit - maker is a one-shot program, so
	// there's no long-lived process to flush on a timer later.
	defer func() {
		if err := shutdown(ctx); err != nil {
			fmt.Println("Error shutting down tracer provider:", err)
		}
	}()

	tracer := otel.Tracer("maker")
	ctx, span := tracer.Start(ctx, "maker.run")
	defer span.End()

	config := config.LoadConfig()
	fmt.Printf("Person Repo URL: %s\n", config.PersonURL)

	repo := people.NewInMemoryPersonRepository(ctx, config.PersonURL)
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
	if err := tracing.PublishWithTrace(ctx, js, "transactions.new", payload); err != nil {
		fmt.Println("Error publishing transaction to NATS:", err)
	} else {
		fmt.Println("Transaction published to NATS")
	}

}

func getRandomAmount() float64 {
	return float64(0 + (uuid.New().ID() % 1000)) // Random amount between 0 and 999
}
