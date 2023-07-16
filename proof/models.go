package proof

import (
	"fmt"
	"strconv"

	"github.com/sirdeggen/merkle/hash"
)

// MerkleProof is the response message for each txid query
type MerkleProofTSC struct {
	TxID       string   `json:"txid"`
	Index      uint64   `json:"index"`
	Root       string   `json:"root"`
	Nodes      []string `json:"nodes"`
	TargetType string   `json:"targetType,omitempty"`
	ProofType  string   `json:"proofType,omitempty"`
	Composite  bool     `json:"composite,omitempty"`
}

type MerkleProof struct {
	Txid      string   `json:"txid"`
	Blockhash string   `json:"blockhash"`
	Index     string   `json:"index"`
	Path      []string `json:"path"`
}

func CreateMerklePathFromBranchesAndIndex(leaves [][]hash.Hash, index uint64) (*MerkleProof, error) {
	var path MerkleProof
	path.Index = string(index)
	levels := uint64(len(leaves)) - 1
	offset := uint64(0)
	mask := uint64(1) << levels
	for level := levels; level <= levels; level-- {
		subIdx := offset
		if index&mask > 0 {
			offset += 1
		} else {
			subIdx += 1
		}
		if level < levels {
			h := leaves[level][subIdx].StringReverse()
			path.Path = append([]string{h}, path.Path...)
		}
		mask = mask >> 1
		offset = offset << 1
	}
	return &path, nil
}

func CheckMerklePathLeadsToRoot(txid *hash.Hash, path *MerkleProof, root *hash.Hash) bool {
	// start with txid
	workingHash := *txid
	lsb, err := strconv.ParseUint(path.Index, 10, 64)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// hash with each path branch
	for _, leaf := range path.Path {
		var digest []byte
		leafBytes := hash.FromStringReverse(leaf)
		// if the least significant bit is 1 then the working hash is on the right
		if lsb&1 > 0 {
			digest = append(leafBytes[:], workingHash[:]...)
		} else {
			digest = append(workingHash[:], leafBytes[:]...)
		}
		workingHash = hash.Sha256Sha256(digest)
		lsb = lsb >> 1
	}
	// check result equality with root
	return workingHash == *root
}
