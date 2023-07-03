package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sirdeggen/merkle/helpers"
	"github.com/sirdeggen/merkle/service"
)

func main() {
	block, err := service.GetBlockFromFile("data/midblock.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("There are ", len(block.Txids), " Transactions in this block.")
	branches, err := service.CalculateMerkleBranches(block)
	if err != nil {
		fmt.Println(err)
	}
	m := helpers.Reverse(block.MerkleRoot)
	fmt.Println("Merkle Root: ", hex.EncodeToString(m[:]))

	cm := helpers.Reverse(branches[len(branches)-1][0])
	fmt.Println("Calculated Merkle Root: ", hex.EncodeToString(cm[:]))
	
}
