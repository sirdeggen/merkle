package merkle

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

type MerkleProofService interface {
	GetMerkleProof(txid string) (*MerkleProof, error)
	StoreMerkleProof(txid string, proof *MerkleProof) error
}