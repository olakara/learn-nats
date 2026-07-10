package main

import (
	"accumulator/account"
	"accumulator/tracing"
	"accumulator/transactions"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type AccountService struct {
	accountRepository account.AccountRepository
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown, err := tracing.Init(ctx, "accumulator")
	if err != nil {
		fmt.Println("Error initializing tracing:", err)
		return
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
		fmt.Println("Error getting JetStream context:", err)
		return
	}

	js.AddConsumer("accumulator", &nats.ConsumerConfig{
		Durable: "transactions-accumulator",
	})

	accountService := &AccountService{
		accountRepository: account.NewInMemoryAccountRepository(),
	}

	js.Subscribe("transactions.new", func(msg *nats.Msg) {
		msgCtx := tracing.ExtractContext(ctx, msg)
		_, span := tracing.StartConsumerSpan(msgCtx, "transactions.new")
		defer span.End()

		var data transactions.Transaction
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
			span.RecordError(err)
		}

		acc, err := accountService.accountRepository.GetByPersonId(data.PersonId)
		if err != nil {
			fmt.Println("No account for person ID:", data.PersonId)
			acc = account.NewAccount(uuid.New().String(), data.PersonId)
			accountService.accountRepository.CreateAccount(acc)
		}
		err = accountService.accountRepository.Accumulate(acc.Id, data.Amount)
		if err != nil {
			fmt.Println("Error accumulating amount:", err)
			span.RecordError(err)
		}

		fmt.Printf("account %+v\n", acc)
	})

	http.Handle("/accounts/{id}", otelhttp.NewHandler(http.HandlerFunc(accountService.getAccountHandler), "get-account"))
	http.Handle("/accounts/", otelhttp.NewHandler(http.HandlerFunc(accountService.getAllAccountsHandler), "get-all-accounts"))

	server := &http.Server{Addr: ":8080"}
	go func() {
		fmt.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error running HTTP server:", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down...")
	server.Shutdown(context.Background())

	if err := shutdown(context.Background()); err != nil {
		fmt.Println("Error shutting down tracer provider:", err)
	}
}

func (s *AccountService) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	account, err := s.accountRepository.GetByPersonId(id)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (s *AccountService) getAllAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.accountRepository.GetAll()
	if err != nil {
		http.Error(w, "Error retrieving accounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
