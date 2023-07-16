package proof

type MerkleProofService interface {
	GetMerkleProof(txid string) (*MerkleProof, error)
	StoreMerkleProof(txid string, proof *MerkleProof) error
}

type merkleProofService struct {
	config string
}

func NewMerkleProofService() *merkleProofService {
	return &merkleProofService{
		config: "no idea waht to put here",
	}
}

func (m *merkleProofService) GetMerkleProof(txids string) (*MerkleProof, error) {
	var proof MerkleProof
	return &proof, nil
}

func (m *merkleProofService) StoreMerkleProof(txid string, proof *MerkleProof) error {
	return nil
}
