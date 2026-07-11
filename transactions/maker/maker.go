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
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type MakerService struct {
	personRepo people.PersonRepository
	js         nats.JetStreamContext
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown, err := tracing.Init(ctx, "maker")
	if err != nil {
		log.Fatal("Error initializing tracing:", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			fmt.Println("Error shutting down tracer provider:", err)
		}
	}()

	cfg := config.LoadConfig()
	fmt.Printf("Person Repo URL: %s\n", cfg.PersonURL)

	repo := people.NewInMemoryPersonRepository(ctx, cfg.PersonURL)
	fmt.Printf("Loaded %d people\n", len(repo.People()))

	if len(repo.People()) == 0 {
		fmt.Println("No people found. Please check the URL or the data source.")
		os.Exit(1)
	}

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

	service := &MakerService{
		personRepo: repo,
		js:         js,
	}

	http.Handle("/start", otelhttp.NewHandler(http.HandlerFunc(service.startHandler), "start-transaction"))

	server := &http.Server{Addr: cfg.Addr}
	go func() {
		fmt.Printf("Server running on %s\n", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error running HTTP server:", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down...")
	server.Shutdown(context.Background())
}

func (s *MakerService) startHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("maker")
	ctx, span := tracer.Start(ctx, "maker.start")
	defer span.End()

	p := s.personRepo.GetSinglePerson()
	t := transaction.NewTransaction(uuid.New().String(), getRandomAmount(), p.Id)
	fmt.Printf("Created Transaction: ID=%s, Amount=%.2f, PersonID=%s\n", t.Id, t.Amount, t.PersonId)

	payload, err := json.Marshal(t)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Error creating transaction", http.StatusInternalServerError)
		return
	}

	if err := tracing.PublishWithTrace(ctx, s.js, "transactions.new", payload); err != nil {
		fmt.Println("Error publishing transaction to NATS:", err)
		span.RecordError(err)
		http.Error(w, "Error publishing transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func getRandomAmount() float64 {
	return float64(0 + (uuid.New().ID() % 1000)) // Random amount between 0 and 999
}
