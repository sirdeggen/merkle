package block

import (
	"math"

	"github.com/sirdeggen/merkle/hash"
)

func BlockBinaryFromJson(blockJson *BlockJson) (*BlockBinary, error) {
	txids := make([]hash.Hash, len(blockJson.Txids))
	for i, hexTxid := range blockJson.Txids {
		txids[i] = hash.FromStringReverse(hexTxid)
	}
	return &BlockBinary{
		Txids:      txids,
		Hash:       hash.FromStringReverse(blockJson.Hash),
		MerkleRoot: hash.FromStringReverse(blockJson.MerkleRoot),
	}, nil
}

func CalculateMerkleBranches(block *BlockBinary) ([][]hash.Hash, error) {
	numberofLevels := int(math.Ceil(math.Log2(float64(len(block.Txids)))))
	// the branches are all 32 byte hashes
	var branches [][]hash.Hash

	// the first layer of branches are just the transactions hashes themselves
	branches = append(branches, block.Txids)

	// the other branches will need to be calculated, and put in these levels of slices
	for i := 0; i < numberofLevels; i++ {
		targetBranches := make([]hash.Hash, 0)
		branches = append(branches, targetBranches)
	}

	// if there's only one then that's the only branch and it's also the root
	if len(block.Txids) == 1 {
		return branches, nil
	}

	for level, branchesAtThisLevel := range branches {
		if len(branchesAtThisLevel) == 1 {
			break
		}
		if level == len(branches)-1 {
			targetBranches := make([]hash.Hash, 0)
			branches = append(branches, targetBranches)
		}
		targetLevel := level + 1
		for idx, branch := range branchesAtThisLevel {
			if idx&1 > 0 {
				digest := append(branchesAtThisLevel[idx-1][:], branch[:]...)
				branches[targetLevel] = append(branches[targetLevel], hash.Sha256Sha256(digest))
				continue
			}
			if idx == len(branchesAtThisLevel)-1 {
				digest := append(branch[:], branch[:]...)
				branches[targetLevel] = append(branches[targetLevel], hash.Sha256Sha256(digest))
				continue
			}
		}
	}
	return branches, nil
}
