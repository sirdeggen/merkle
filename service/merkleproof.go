package main

import (
	merkle "github.com/sirdeggen/merkle/models"
)

type merkleProofService struct {
	config string
}

func NewMekleProofService() *merkleProofService {
	return &merkleProofService{
		config: "no idea waht to put here",
	}
}

func (m *merkleProofService) GetMerkleProof(txids string) (*merkle.MerkleProof, error) {
	var proof merkle.MerkleProof
	return &proof, nil
}

func (m *merkleProofService) StoreMerkleProof(txid string, proof *merkle.MerkleProof) error {
	return nil
}
