package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirdeggen/merkle/tree"
)

func readExisting(root string) {
	mts := tree.NewMerkleTreeService("data")
	x := uint64(0)
	for {
		path, err := mts.Read(root, x)
		if err != nil {
			fmt.Println(err)
			break
		}
		jsonPath := path.Json()
		stringBytes, err := json.MarshalIndent(jsonPath, "", "    ")
		fmt.Println(string(stringBytes))
		x++
	}
}

func createTreeFileFromJsonFile(filename string) {
	mts := tree.NewMerkleTreeService("data")
	var branches tree.MerkleTree
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &branches)
	if err != nil {
		fmt.Println(err)
	}
	err = mts.Write(branches)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	readExisting("d66e56fb408763e36e8622eb56a8a1072ccc606476fe9e0765cca0dff95949b1")
	createTreeFileFromJsonFile("data/smallblock.json")
}
