package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var nblock int

type Block struct {
	From         string
	To           string
	Amount       int
	PreviousHash string
}

func (b *Block) toString() string {

	return b.From + " -> " + b.To + " : " + strconv.Itoa(b.Amount) + "\n" + b.PreviousHash
}

func (b *Block) save() error {
	data := []byte(b.toString())
	fname := "blocks/" + strconv.Itoa(nblock) + ".block"

	fmt.Println(fname)

	if err := ioutil.WriteFile(fname, data, 0600); err != nil {
		return err
	}

	return nil
}

// Hash block # n
func hash(n int) (string, error) {
	fname := "blocks/" + strconv.Itoa(n) + ".block"
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}

	sha := sha256.New()
	_, err = sha.Write(f)
	if err != nil {
		return "", err
	}

	fmt.Printf("%x", sha.Sum(nil))

	h := fmt.Sprintf("%x", sha.Sum(nil))

	return h, nil
}

// Read n numbered block
func read(n int) (*Block, error) {
	fname := "blocks/" + strconv.Itoa(n) + ".block"
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Scan the first line: transaction
	scanner.Scan()
	firstLine := strings.Split(scanner.Text(), " ")

	amount, err := strconv.ParseInt(firstLine[4], 10, 32)
	if err != nil {
		return nil, err
	}

	// Scan the next line: previous hash
	scanner.Scan()
	nextLine := scanner.Text()

	fmt.Println("next line\t" + nextLine)

	b := Block{From: firstLine[0], To: firstLine[2], Amount: int(amount), PreviousHash: nextLine}

	return &b, nil
}

// Calculate last transaction hash
// Fill in the transaction
// Add the last hash
// Save
func MakeTransaction(from, to string, amount int) error {
	lh, err := hash(nblock - 1)
	if err != nil {
		return err
	}

	b := &Block{From: from, To: to, Amount: amount, PreviousHash: lh}
	if err := b.save(); err != nil {
		return err
	}

	nblock++

	return nil
}

func ValidateChain(nblock int) bool {
	// calculate hash of nblock-1, compare it to hash from n block
	// do until genezis (0)
	for nblock > 0 {
		h, err := hash(nblock - 1)
		if err != nil {
			panic(err)
		}

		b, err := read(nblock)
		if err != nil {
			panic(err)
		}

		if b.PreviousHash != h {
			return false
		}

		nblock--
	}
	return true
}

func main() {
	nblock = 1

	if err := MakeTransaction("Mike", "Erika", 15); err != nil {
		panic(err)
	}
	if err := MakeTransaction("Steve", "John", 2); err != nil {
		panic(err)
	}
	if err := MakeTransaction("Morr", "Jom", 9); err != nil {
		panic(err)
	}

	fmt.Print(ValidateChain(nblock - 1))
}
