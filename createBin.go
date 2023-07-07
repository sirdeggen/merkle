package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"

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

	var fileBytes []byte
	numOfTxs := make([]byte, 8)
	binary.LittleEndian.PutUint64(numOfTxs, uint64(len(block.Txids)))

	fileBytes = append(fileBytes, numOfTxs...)
	for x := len(branches) - 2; x >= 0; x-- {
		for y := 0; y < len(branches[x]); y++ {
			fileBytes = append(fileBytes, branches[x][y][:]...)
		}
	}
	// write to file
	err = ioutil.WriteFile("data/branches.bin", fileBytes, fs.FileMode(0644))
	if err != nil {
		fmt.Println(err)
	}
}
