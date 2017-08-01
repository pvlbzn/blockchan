package main

import (
	"fmt"

	"github.com/pvlbzn/blockchan/block"
)

func main() {

	for n := 0; n < 5; n++ {
		b, _ := block.NewBlock()
		b.AddTransaction(block.NewTransaction("Sam", "Moe", n*3))
		b.AddTransaction(block.NewTransaction("Carter", "Jimmy", n*5))
		b.AddTransaction(block.NewTransaction("Merlin", "Foe", n*2))
		b.AddTransaction(block.NewTransaction("Lou", "Joe", n*9))
		b.Write()
	}

	s, _ := block.ValidateChain()
	fmt.Println(s)
}
