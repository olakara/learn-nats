package account

import "fmt"

type Account struct {
	Id       string
	Balance  float64
	PersonId string
}

func NewAccount(id string, personId string) *Account {
	return &Account{
		Id:       id,
		Balance:  0,
		PersonId: personId,
	}
}

type AccountRepository interface {
	CreateAccount(account *Account) error
	Accumulate(accountId string, amount float64) error
	GetByPersonId(personId string) (*Account, error)
	GetAll() ([]*Account, error)
}

type InMemoryAccountRepository struct {
	accounts map[string]*Account
}

func (r *InMemoryAccountRepository) CreateAccount(account *Account) error {
	if _, exists := r.accounts[account.Id]; exists {
		return fmt.Errorf("account already exists")
	}
	r.accounts[account.Id] = account
	return nil
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[string]*Account),
	}
}

func (r *InMemoryAccountRepository) Accumulate(accountId string, amount float64) error {
	account, exists := r.accounts[accountId]
	if !exists {
		return fmt.Errorf("account not found")
	}
	account.Balance += amount
	return nil
}

func (r *InMemoryAccountRepository) GetByPersonId(personId string) (*Account, error) {
	for _, account := range r.accounts {
		if account.PersonId == personId {
			return account, nil
		}
	}
	return nil, fmt.Errorf("account not found")
}

func (r *InMemoryAccountRepository) GetAll() ([]*Account, error) {
	accounts := make([]*Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}
