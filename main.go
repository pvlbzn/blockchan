package main

import (
	"fmt"

	"github.com/pvlbzn/blockchan/block"
)

func main() {
	b, _ := block.NewBlock()
	b.AddTransaction(block.NewTransaction("Sam", "Moe", 51))
	b.Write()

	for n := 0; n < 5; n++ {
		b, _ := block.NewBlock()
		b.AddTransaction(block.NewTransaction("Sam", "Moe", n*3))
		b.Write()
	}

	s, _ := block.ValidateChain()
	fmt.Println(s)
}
