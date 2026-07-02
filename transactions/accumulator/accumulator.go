package main

import (
	"accumulator/account"
	"accumulator/transactions"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type AccountService struct {
	accountRepository account.AccountRepository
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
		var data transactions.Transaction
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			fmt.Println("Error unmarshalling data:", err)
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
		}

		fmt.Printf("account %+v\n", acc)
	})

	http.HandleFunc("/accounts/{id}", accountService.getAccountHandler)
	http.HandleFunc("/accounts/", accountService.getAllAccountsHandler)

	// Start the server on port 8080
	println("Server running on :8080")
	http.ListenAndServe(":8080", nil)

	select {}

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
