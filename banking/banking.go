package banking

import (
	"errors"
	"fmt"
)

// Account Struct
type Account struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("NO MONEY")

// Create constructor
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// dont copy!!!!! in receiver, use real account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

func (a *Account) Withraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

func (a Account) Balance() int {
	return a.balance
}

// Change owner
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// return owner
func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return fmt.Sprint(a.owner, "'s account.\nHAs: ", a.balance)
}
