package main

import (
	"fmt"
	"maker/config"
	"maker/people"
	"maker/transaction"
	"os"

	"github.com/google/uuid"
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

}
