package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirdeggen/merkle/block"
	"github.com/sirdeggen/merkle/tree"
)

func readExisting(root string, index int) {
	mts := tree.NewMerkleTreeService("data")

	x := uint64(index)
	path, err := mts.Read(root, x)
	if err != nil {
		fmt.Println(err)
	}
	jsonPath := path.Json()
	stringBytes, err := json.MarshalIndent(jsonPath, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(stringBytes))
}

func createTreeFileFromJsonFile(filename string) {
	file, _ := os.Open(filename)
	defer file.Close()

	var jsonData block.BlockJson
	byteValue, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(byteValue, &jsonData)

	blockBinary, _ := block.BlockBinaryFromJson(&jsonData)

	branches, _ := blockBinary.CalculateMerkleBranches()
	mts := tree.NewMerkleTreeService("data")
	err := mts.Write(branches)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	createTreeFileFromJsonFile("data/smallblock.json")
	createTreeFileFromJsonFile("data/midblock.json")
	createTreeFileFromJsonFile("data/block.json")
	readExisting("0a3b8cb97063d49e1a1504f10c5c6e648ec8fc436f1ab0ee68dd457a305f0a8b", 1)    // small
	readExisting("d66e56fb408763e36e8622eb56a8a1072ccc606476fe9e0765cca0dff95949b1", 12)   // med
	readExisting("a623039e2030dfafd02af3948f8f8483aadbb7296e205ff95f78a52269be97f5", 1234) // big
}
