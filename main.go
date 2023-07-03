package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/sirdeggen/merkle/helpers"
	"github.com/sirdeggen/merkle/models"
	"github.com/sirdeggen/merkle/service"
)

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

func getBlockFromFile(filename string) (*models.BlockBinary, error) {
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

func calculateMerkleNodes(block *models.BlockBinary) ([][]models.Hash, error) {
	numberofLevels := int(math.Ceil(math.Log2(float64(len(block.Txids)))))
	// the nodes are all 32 byte hashes
	var nodes [][]models.Hash

	// the first layer of nodes are just the transactions hashes themselves
	nodes = append(nodes, block.Txids)

	// the other nodes will need to be calculated, and put in these levels of slices
	for i := 0; i < numberofLevels; i++ {
		targetNodes := make([]models.Hash, 0)
		nodes = append(nodes, targetNodes)
	}

	// if there's only one then that's the only node and it's also the root
	if len(block.Txids) == 1 {
		return nodes, nil
	}

	for level, nodesAtThisLevel := range nodes {
		if len(nodesAtThisLevel) == 1 {
			break
		}
		if level == len(nodes)-1 {
			targetNodes := make([]models.Hash, 0)
			nodes = append(nodes, targetNodes)
		}
		targetLevel := level + 1
		for idx, node := range nodesAtThisLevel {
			if idx&1 > 0 {
				digest := append(nodesAtThisLevel[idx-1][:], node[:]...)
				nodes[targetLevel] = append(nodes[targetLevel], H(digest))
				continue
			}
			if idx == len(nodesAtThisLevel)-1 {
				digest := append(node[:], node[:]...)
				nodes[targetLevel] = append(nodes[targetLevel], H(digest))
				continue
			}
		}
		numOfNodesAtTargetLevel := len(nodes[targetLevel])
		if numOfNodesAtTargetLevel > 1 && numOfNodesAtTargetLevel&1 > 0 {
			nodes[targetLevel] = append(nodes[targetLevel], nodes[targetLevel][len(nodes[targetLevel])-1])
		}
	}
	return nodes, nil
}

func createMerklePathFromNodesAndIndex(nodes [][]models.Hash, index uint64) (*models.MerklePathBinary, error) {
	var path models.MerklePathBinary
	path.Index = index
	levels := uint64(len(nodes)) - 1
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
			path.Nodes = append([]models.Hash{nodes[level][subIdx]}, path.Nodes...)
		}
		mask = mask >> 1
		offset = offset << 1
	}
	return &path, nil
}

func checkMerklePathLeadsToRoot(txid *models.Hash, path *models.MerklePathBinary, root *models.Hash) bool {
	// start with txid
	workingHash := *txid
	lsb := path.Index
	// hash with each path node
	for _, node := range path.Nodes {
		var digest []byte
		// if the least significant bit is 1 then the working hash is on the right
		if lsb&1 > 0 {
			digest = append(node[:], workingHash[:]...)
		} else {
			digest = append(workingHash[:], node[:]...)
		}
		workingHash = H(digest)
		lsb = lsb >> 1
	}
	// check result equality with root
	return workingHash == *root
}

func calculateBlockWideMerklePaths(block *models.BlockBinary) error {
	nodes, err := calculateMerkleNodes(block)
	if err != nil {
		return err
	}
	pathmap := make(models.PathMap)
	for idx, txid := range block.Txids {
		path, err := createMerklePathFromNodesAndIndex(nodes, uint64(idx))
		if err != nil {
			fmt.Println(err)
		}
		pathmap[txid] = path
	}
	block.MerklePaths = &pathmap
	return nil
}

func jsonBlockFromBinary(block *models.BlockBinary) (*models.BlockJson, error) {
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
		for _, node := range path.Nodes {
			revNode := helpers.Reverse(node)
			mpJ.Nodes = append(mpJ.Nodes, hex.EncodeToString(revNode[:]))
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

func main() {
	serv := service.NewMerkleProofService()
	fmt.Println(serv)

	block, err := getBlockFromFile("data/block.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("There are ", len(block.Txids), " Transactions in this block.")

	nodes, err := calculateMerkleNodes(block)
	if err != nil {
		fmt.Println(err)
	}
	m := helpers.Reverse(block.MerkleRoot)
	fmt.Println("Merkle Root: ", hex.EncodeToString(m[:]))

	cm := helpers.Reverse(nodes[len(nodes)-1][0])
	fmt.Println("Calculated Merkle Root: ", hex.EncodeToString(cm[:]))

	// create the merkle path
	path, err := createMerklePathFromNodesAndIndex(nodes, 11)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Merkle Path: ")
	for _, node := range path.Nodes {
		rev := helpers.Reverse(node)
		fmt.Println(hex.EncodeToString(rev[:]))
	}

	root := block.MerkleRoot
	txid := block.Txids[path.Index]
	valid := checkMerklePathLeadsToRoot(&txid, path, &root)
	fmt.Println("Merkle Proof Valid: ", valid)

	fmt.Println("Calculating block wide merkle paths...")
	err = calculateBlockWideMerklePaths(block)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done.")
	jBlock, err := jsonBlockFromBinary(block)
	if err != nil {
		fmt.Println(err)
	}
	jsonString, err := json.MarshalIndent(jBlock, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(*block.MerklePaths))
	fmt.Println(string(jsonString))
}
