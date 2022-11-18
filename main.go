package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type BlockJson struct {
	Txids      []string `json:"tx"`
	Hash       string   `json:"hash"`
	MerkleRoot string   `json:"merkleroot"`
}

type BlockBinary struct {
	Txids      [][]byte
	Hash       []byte
	MerkleRoot []byte
}

func hexToBytes(h string) ([]byte, error) {
	bytes, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func blockBinaryFromJson(blockJson *BlockJson) (*BlockBinary, error) {
	txids := make([][]byte, len(blockJson.Txids))
	for i, txid := range blockJson.Txids {
		txids[i] = []byte(txid)
	}
	hash, err := hexToBytes(blockJson.Hash)
	if err != nil {
		return nil, err
	}
	merkleRoot, err := hexToBytes(blockJson.MerkleRoot)
	if err != nil {
		return nil, err
	}
	return &BlockBinary{
		Txids:      txids,
		Hash:       hash,
		MerkleRoot: merkleRoot,
	}, nil
}

func getBlockFromFile(filename string) (*BlockBinary, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var jsonData BlockJson
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}
	opti, err := blockBinaryFromJson(&jsonData)
	if err != nil {
		return nil, err
	}
	return opti, nil
}

func main() {
	fmt.Println("Hello, World!")
	block, err := getBlockFromFile("data/block.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(block)
}
