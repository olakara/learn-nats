package transactions

import "time"

type Transaction struct {
	Id        string    `json:"id"`
	Amount    float64   `json:"amount"`
	PersonId  string    `json:"personId"`
	Timestamp time.Time `json:"timestamp"`
}
