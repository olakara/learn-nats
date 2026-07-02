package transaction

type Transaction struct {
	Id       string
	Amount   float64
	PersonId string
}

func NewTransaction(id string, amount float64, personId string) Transaction {
	return Transaction{Id: id, Amount: amount, PersonId: personId}
}
