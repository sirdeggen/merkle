package models

// MerkleProof is the response message for each txid query
type MerkleProof struct {
	TxID       string   `json:"txid"`
	Index      uint64   `json:"index"`
	Root       string   `json:"root"`
	Nodes      []string `json:"nodes"`
	TargetType string   `json:"targetType,omitempty"`
	ProofType  string   `json:"proofType,omitempty"`
	Composite  bool     `json:"composite,omitempty"`
}

// MerklePath optimized for responding to requests quickly
type MerklePath struct {
	Nodes []*Hash `json:"nodes"`
	Index uint64  `json:"index"`
}

// MerkleBlock optimized for storing the necessary data
type MerkleBlock struct {
	Hash       Hash
	Root       Hash
	MerkleTree []Hash
	PathMap    map[Hash]MerklePath
}

type MerklePathBinary struct {
	Nodes []Hash `json:"nodes"`
	Index uint64 `json:"index"`
}

type MerklePathJson struct {
	Nodes []string `json:"nodes"`
	Index uint64   `json:"index"`
}

type PathMap map[Hash]*MerklePathBinary

type MerkleProofService interface {
	GetMerkleProof(txid string) (*MerkleProof, error)
	StoreMerkleProof(txid string, proof *MerkleProof) error
}
