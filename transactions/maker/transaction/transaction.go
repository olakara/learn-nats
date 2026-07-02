package transaction

import "time"

type Transaction struct {
	Id        string    `json:"id"`
	Amount    float64   `json:"amount"`
	PersonId  string    `json:"personId"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTransaction(id string, amount float64, personId string) Transaction {
	return Transaction{Id: id, Amount: amount, PersonId: personId, Timestamp: time.Now()}
}
