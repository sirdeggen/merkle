package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/sirdeggen/merkle/helpers"
	"github.com/sirdeggen/merkle/service"
)

func main() {
	serv := service.NewMerkleProofService()
	fmt.Println(serv)

	block, err := service.GetBlockFromFile("data/block.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("There are ", len(block.Txids), " Transactions in this block.")

	nodes, err := service.CalculateMerkleNodes(block)
	if err != nil {
		fmt.Println(err)
	}
	m := helpers.Reverse(block.MerkleRoot)
	fmt.Println("Merkle Root: ", hex.EncodeToString(m[:]))

	cm := helpers.Reverse(nodes[len(nodes)-1][0])
	fmt.Println("Calculated Merkle Root: ", hex.EncodeToString(cm[:]))

	// create the merkle path
	path, err := service.CreateMerklePathFromNodesAndIndex(nodes, 11)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Merkle Path: ")
	for _, node := range path.Nodes {
		rev := helpers.Reverse(node)
		fmt.Println(hex.EncodeToString(rev[:]))
	}

	root := block.MerkleRoot
	txid := block.Txids[path.Index]
	valid := service.CheckMerklePathLeadsToRoot(&txid, path, &root)
	fmt.Println("Merkle Proof Valid: ", valid)

	fmt.Println("Calculating block wide merkle paths...")
	err = service.CalculateBlockWideMerklePaths(block)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done.")
	jBlock, err := service.JsonBlockFromBinary(block)
	if err != nil {
		fmt.Println(err)
	}
	jsonString, err := json.MarshalIndent(jBlock, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(*block.MerklePaths))
	fmt.Println(string(jsonString))
}
