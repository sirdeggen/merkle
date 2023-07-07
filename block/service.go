package block

import "github.com/sirdeggen/merkle/hash"

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
