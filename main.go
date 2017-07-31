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
	return b.From + " -> " + b.To + " : " + strconv.Itoa(b.Amount)
}

func (b *Block) save() error {
	data := []byte(b.toString())
	fname := "blocks/" + strconv.Itoa(nblock) + ".block"

	if err := ioutil.WriteFile(fname, data, 0600); err != nil {
		return err
	}

	return nil
}

// Hash block # n
func hash(n int) error {
	fname := "blocks/" + strconv.Itoa(n) + ".block"
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	sha := sha256.New()
	_, err = sha.Write(f)
	if err != nil {
		return err
	}

	fmt.Printf("%x", sha.Sum(nil))

	return nil
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

	b := Block{From: firstLine[0], To: firstLine[2], Amount: int(amount), PreviousHash: nextLine}

	return &b, nil
}

func main() {
	nblock = 1

	// b := Block{From: "Chan", To: "Cat", Amount: 50}

	// b.save()

	bn, _ := read(0)
	fmt.Println(bn.toString())
	hash(0)
}
