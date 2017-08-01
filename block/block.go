package block

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Block struct {
	Transactions []Transaction
	PreviousHash string
	Number       int
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

// String implements Stringer interface. Such that a block
// can be printed well using fmt package.
func (b *Block) String() string {
	var t string
	for _, transaction := range b.Transactions {
		t += fmt.Sprintf("%s\n", transaction.String())
	}
	return t + b.PreviousHash
}

// NewTransaction is a convinience wrapper for a new transaction.
func NewTransaction(from, to string, amount int) *Transaction {
	return &Transaction{From: from, To: to, Amount: amount}
}

// NewBlock constructor for a block. It calculates a hash of a last block
// and finds a next number of a block.
func NewBlock() (*Block, error) {
	n, err := findLast()
	if err != nil {
		return nil, err
	}

	h, err := ftoh(n)
	if err != nil {
		return nil, err
	}

	fmt.Println("Hash: " + h)

	return &Block{Number: n + 1, PreviousHash: h}, nil
}

// findLast block in blocks directory
func findLast() (int, error) {
	blocks, err := ioutil.ReadDir("blocks")
	if err != nil {
		return -1, err
	}

	var n int
	for _, b := range blocks {
		bn := strings.Split(b.Name(), ".")[0]
		number, err := strconv.Atoi(bn)
		if err != nil {
			continue
		}

		if number > n {
			n = number
		}
	}

	return n, nil
}

// ftoh (file to hash)
func ftoh(n int) (string, error) {
	f, err := ioutil.ReadFile("blocks/" + strconv.Itoa(n) + ".block")
	if err != nil {
		return "", err
	}

	return hash(f)
}

// hash bytes -> hash string
func hash(f []byte) (string, error) {
	sha := sha256.New()
	_, err := sha.Write(f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha.Sum(nil)), nil
}

// AddTransaction to transaction collection
func (b *Block) AddTransaction(t *Transaction) {
	b.Transactions = append(b.Transactions, *t)
}

// Write a block into a file
func (b *Block) Write() error {
	data := []byte(b.String())
	fname := "blocks/" + strconv.Itoa(b.Number) + ".block"

	if err := ioutil.WriteFile(fname, data, 0600); err != nil {
		return err
	}

	return nil
}

// Read a block from a file
func Read(n int) (*Block, error) {
	fname := "blocks/" + strconv.Itoa(n) + ".block"
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Scan each line into data array
	scanner := bufio.NewScanner(f)
	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	// Read hash
	// Hash is known to be the last line in the block
	h := data[len(data)-1]

	// Read transactions
	var transactions []Transaction
	for n := len(data) - 2; n > 0; n-- {
		l := strings.Split(data[n], " ")
		amount, err := strconv.ParseInt(l[4], 10, 32)
		if err != nil {
			return nil, err
		}
		t := NewTransaction(l[0], l[2], int(amount))
		transactions = append(transactions, *t)
	}

	return &Block{Transactions: transactions, PreviousHash: h, Number: n}, nil
}

// ValidateChain checks chain hash validity
func ValidateChain() (bool, error) {
	n, err := findLast()
	if err != nil {
		return false, err
	}

	for n > 0 {
		h, err := ftoh(n - 1)
		if err != nil {
			return false, err
		}

		b, err := Read(n)
		if err != nil {
			return false, err
		}

		if b.PreviousHash != h {
			return false, nil
		}

		n--
	}

	return true, nil
}
