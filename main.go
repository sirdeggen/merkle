package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirdeggen/merkle/block"
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
	file, _ := os.Open(filename)
	defer file.Close()

	var jsonData block.BlockJson
	byteValue, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(byteValue, &jsonData)

	blockBinary, _ := block.BlockBinaryFromJson(&jsonData)

	fmt.Println("BlockBinary:", blockBinary)
}

func main() {
	readExisting("d66e56fb408763e36e8622eb56a8a1072ccc606476fe9e0765cca0dff95949b1")
	createTreeFileFromJsonFile("data/smallblock.json")
}
