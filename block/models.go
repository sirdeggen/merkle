package block

import "github.com/sirdeggen/merkle/hash"

type BlockJson struct {
	Txids      []string `json:"tx"`
	Hash       string   `json:"hash"`
	MerkleRoot string   `json:"merkleroot"`
}

type BlockBinary struct {
	Txids      []hash.Hash
	Hash       hash.Hash
	MerkleRoot hash.Hash
}
