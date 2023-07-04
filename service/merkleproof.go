package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/sirdeggen/merkle/helpers"
	"github.com/sirdeggen/merkle/models"
)

type merkleProofService struct {
	config string
}

func NewMerkleProofService() *merkleProofService {
	return &merkleProofService{
		config: "no idea waht to put here",
	}
}

func (m *merkleProofService) GetMerkleProof(txids string) (*models.MerkleProof, error) {
	var proof models.MerkleProof
	return &proof, nil
}

func (m *merkleProofService) StoreMerkleProof(txid string, proof *models.MerkleProof) error {
	return nil
}

func blockBinaryFromJson(blockJson *models.BlockJson) (*models.BlockBinary, error) {
	txids := make([]models.Hash, len(blockJson.Txids))
	for i, hexTxid := range blockJson.Txids {
		var txid [32]byte
		hash, err := helpers.HexToBytes(hexTxid)
		if err != nil {
			return nil, err
		}
		copy(txid[:], []byte(hash))
		txids[i] = helpers.Reverse(txid)
	}
	hash, err := helpers.HexToBytes(blockJson.Hash)
	if err != nil {
		return nil, err
	}
	merkleRoot, err := helpers.HexToBytes(blockJson.MerkleRoot)
	if err != nil {
		return nil, err
	}
	var m [32]byte
	copy(m[:], []byte(merkleRoot))
	return &models.BlockBinary{
		Txids:      txids,
		Hash:       hash,
		MerkleRoot: helpers.Reverse(m),
	}, nil
}

func H(digest []byte) [32]byte {
	one := sha256.Sum256(digest)
	return sha256.Sum256(one[:])
}

func GetBlockFromFile(filename string) (*models.BlockBinary, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var jsonData models.BlockJson
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}
	opti, err := blockBinaryFromJson(&jsonData)
	if err != nil {
		return nil, err
	}
	return opti, nil
}

func CalculateMerkleBranches(block *models.BlockBinary) ([][]models.Hash, error) {
	numberofLevels := int(math.Ceil(math.Log2(float64(len(block.Txids)))))
	// the branches are all 32 byte hashes
	var branches [][]models.Hash

	// the first layer of branches are just the transactions hashes themselves
	branches = append(branches, block.Txids)

	// the other branches will need to be calculated, and put in these levels of slices
	for i := 0; i < numberofLevels; i++ {
		targetBranches := make([]models.Hash, 0)
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
			targetBranches := make([]models.Hash, 0)
			branches = append(branches, targetBranches)
		}
		targetLevel := level + 1
		for idx, branch := range branchesAtThisLevel {
			if idx&1 > 0 {
				digest := append(branchesAtThisLevel[idx-1][:], branch[:]...)
				branches[targetLevel] = append(branches[targetLevel], H(digest))
				continue
			}
			if idx == len(branchesAtThisLevel)-1 {
				digest := append(branch[:], branch[:]...)
				branches[targetLevel] = append(branches[targetLevel], H(digest))
				continue
			}
		}
		// numOfBranchesAtTargetLevel := len(branches[targetLevel])
		// if numOfBranchesAtTargetLevel > 1 && numOfBranchesAtTargetLevel&1 > 0 {
		// 	branches[targetLevel] = append(branches[targetLevel], branches[targetLevel][len(branches[targetLevel])-1])
		// }
	}
	return branches, nil
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

func CheckMerklePathLeadsToRoot(txid *models.Hash, path *models.MerklePathBinary, root *models.Hash) bool {
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
		workingHash = H(digest)
		lsb = lsb >> 1
	}
	// check result equality with root
	return workingHash == *root
}

func CalculateBlockWideMerklePaths(block *models.BlockBinary) error {
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

func JsonBlockFromBinary(block *models.BlockBinary) (*models.BlockJson, error) {
	txids := make([]string, len(block.Txids))
	for i, txid := range block.Txids {
		rev := helpers.Reverse(txid)
		txids[i] = hex.EncodeToString(rev[:])
	}
	var jBlock models.BlockJson
	jBlock.MerklePaths = make(map[string]models.MerklePathJson)
	for txid, path := range *block.MerklePaths {
		var mpJ models.MerklePathJson
		rev := helpers.Reverse(txid)
		for _, leaf := range path.Path {
			revLeaf := helpers.Reverse(leaf)
			mpJ.Path = append(mpJ.Path, hex.EncodeToString(revLeaf[:]))
		}
		mpJ.Index = path.Index
		jBlock.MerklePaths[hex.EncodeToString(rev[:])] = mpJ
	}
	rMerkleRoot := helpers.Reverse(block.MerkleRoot)
	jBlock.Txids = txids
	jBlock.Hash = hex.EncodeToString(block.Hash)
	jBlock.MerkleRoot = hex.EncodeToString(rMerkleRoot[:])
	return &jBlock, nil
}
