package main

import (
	"fmt"

	"github.com/sirdeggen/merkle/tree"
)

func main() {
	mts := tree.NewMerkleTreeService()
	path, err := mts.Read("data/branches.bin", 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*path)
}
