package models

type BlockJson struct {
	Txids       []string                  `json:"tx"`
	Hash        string                    `json:"hash"`
	MerkleRoot  string                    `json:"merkleroot"`
	MerklePaths map[string]MerklePathJson `json:"merklepaths"`
}

type BlockBinary struct {
	Txids       []Hash
	Hash        []byte
	MerkleRoot  Hash
	MerklePaths *PathMap
}
