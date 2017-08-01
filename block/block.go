package block

import (
	"fmt"
	"strconv"
)

type Block struct {
	Transactions []Transaction
	PreviousHash string
	Number       string
}

type Transaction struct {
	From   string
	To     string
	Amount int
}

// String implements Stringer interface. Such that a transaction
// can be printed well using fmt package.
func (t *Transaction) String() string {
	return t.From + " -> " + t.To + " : " + strconv.Itoa(t.Amount)
}

func (b *Block) String() string {
	var t string
	for i, transaction := range b.Transactions {
		t += fmt.Sprintf("%d: %s\n", i+1, transaction.String())
	}
	return t + b.PreviousHash
}

func (b *Block) AddTransaction(t *Transaction) error {

}

func (b *Block) Write() error {

}

func (b *Block) Read() (*Block, error) {

}

func ValidateChain(lastBlock int) bool {

}
