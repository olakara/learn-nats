package people

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"time"
)

type Person struct {
	Id   string
	Name string
}

func NewPerson(id string, name string) Person {
	return Person{Id: id, Name: name}
}

type PersonRepository interface {
	People() []Person
}

type InMemoryPersonRepository struct {
	people []Person
}

func NewInMemoryPersonRepository(url string) *InMemoryPersonRepository {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return &InMemoryPersonRepository{people: []Person{}}
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &InMemoryPersonRepository{people: []Person{}}
	}

	type personDTO struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var data []personDTO
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return &InMemoryPersonRepository{people: []Person{}}
	}

	people := make([]Person, 0, len(data))
	for _, p := range data {
		people = append(people, NewPerson(p.ID, p.Name))
	}

	return &InMemoryPersonRepository{people: people}
}

func (r *InMemoryPersonRepository) People() []Person {
	return r.people
}

func (r *InMemoryPersonRepository) GetSinglePerson() Person {
	if len(r.people) == 0 {
		return Person{}
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(r.people))))
	if err != nil {
		return r.people[0]
	}

	return r.people[int(n.Int64())]
}
