package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

type BlockJson struct {
	Txids      []string `json:"tx"`
	Hash       string   `json:"hash"`
	MerkleRoot string   `json:"merkleroot"`
}

type BlockBinary struct {
	Txids      [][32]byte
	Hash       []byte
	MerkleRoot [32]byte
}

func hexToBytes(h string) ([]byte, error) {
	bytes, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func reverse(b [32]byte) [32]byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func blockBinaryFromJson(blockJson *BlockJson) (*BlockBinary, error) {
	txids := make([][32]byte, len(blockJson.Txids))
	for i, hexTxid := range blockJson.Txids {
		var txid [32]byte
		hash, err := hexToBytes(hexTxid)
		if err != nil {
			return nil, err
		}
		copy(txid[:], []byte(hash))
		txids[i] = reverse(txid)
	}
	hash, err := hexToBytes(blockJson.Hash)
	if err != nil {
		return nil, err
	}
	merkleRoot, err := hexToBytes(blockJson.MerkleRoot)
	if err != nil {
		return nil, err
	}
	var m [32]byte
	copy(m[:], []byte(merkleRoot))
	return &BlockBinary{
		Txids:      txids,
		Hash:       hash,
		MerkleRoot: reverse(m),
	}, nil
}

func H(digest []byte) [32]byte {
	one := sha256.Sum256(digest)
	return sha256.Sum256(one[:])
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

func calculateMerkleNodes(block *BlockBinary) ([][][32]byte, error) {
	numberofLevels := int(math.Ceil(math.Log2(float64(len(block.Txids)))))
	// the nodes are all 32 byte hashes
	var nodes [][][32]byte

	// the first layer of nodes are just the transactions hashes themselves
	nodes = append(nodes, block.Txids)
	for i := 0; i < numberofLevels; i++ {
		targetNodes := make([][32]byte, 0)
		nodes = append(nodes, targetNodes)
	}

	fmt.Println(numberofLevels)
	// fmt.Println(nodes)

	// if there's only one then that's the only node and it's also the root
	if len(block.Txids) == 1 {
		return nodes, nil
	}

	for level, nodesAtThisLevel := range nodes {
		if len(nodesAtThisLevel) == 1 {
			break
		}
		if level == len(nodes)-1 {
			fmt.Println("adding level")
			targetNodes := make([][32]byte, 0)
			nodes = append(nodes, targetNodes)
		}
		targetLevel := level + 1
		var visualization string
		for idx, node := range nodesAtThisLevel {
			visualization += "-"
			if (idx % 2) != 0 {
				digest := append(nodesAtThisLevel[idx-1][:], node[:]...)
				nodes[targetLevel] = append(nodes[targetLevel], H(digest))
				continue
			}
			if idx == len(nodesAtThisLevel)-1 {
				digest := append(node[:], node[:]...)
				nodes[targetLevel] = append(nodes[targetLevel], H(digest))
				continue
			}
		}
		fmt.Println(visualization)
	}
	return nodes, nil
}

func main() {
	block, err := getBlockFromFile("data/block.json")
	if err != nil {
		fmt.Println(err)
	}

	nodes, err := calculateMerkleNodes(block)
	if err != nil {
		fmt.Println(err)
	}
	m := reverse(block.MerkleRoot)
	fmt.Println("Merkle Root: ", hex.EncodeToString(m[:]))

	cm := reverse(nodes[len(nodes)-1][0])
	fmt.Println("Calculated Merkle Root: ", hex.EncodeToString(cm[:]))
}
