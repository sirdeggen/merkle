package proof

import (
	"fmt"

	"github.com/sirdeggen/merkle/block"
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

func CreateMerklePathFromBranchesAndIndex(leaves [][]models.Hash, index uint64) (*models.MerklePathBinary, error) {
	var path models.MerklePathBinary
	path.Index = index
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
			path.Path = append([]models.Hash{leaves[level][subIdx]}, path.Path...)
		}
		mask = mask >> 1
		offset = offset << 1
	}
	return &path, nil
}

func CheckMerklePathLeadsToRoot(txid *hash.Hash, path *block.MerklePathBinary, root *hash.Hash) bool {
	// start with txid
	workingHash := *txid
	lsb := path.Index
	// hash with each path branch
	for _, leaf := range path.Path {
		var digest []byte
		// if the least significant bit is 1 then the working hash is on the right
		if lsb&1 > 0 {
			digest = append(leaf[:], workingHash[:]...)
		} else {
			digest = append(workingHash[:], leaf[:]...)
		}
		workingHash = hash.Sha256Sha256(digest)
		lsb = lsb >> 1
	}
	// check result equality with root
	return workingHash == *root
}

func CalculateBlockWideMerklePaths(block *block.BlockBinary) error {
	branches, err := CalculateMerkleBranches(block)
	if err != nil {
		return err
	}
	pathmap := make(models.PathMap)
	for idx, txid := range block.Txids {
		path, err := CreateMerklePathFromBranchesAndIndex(branches, uint64(idx))
		if err != nil {
			fmt.Println(err)
		}
		pathmap[txid] = path
	}
	block.MerklePaths = &pathmap
	return nil
}
