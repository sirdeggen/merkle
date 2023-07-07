package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"os"

	"github.com/sirdeggen/merkle/helpers"
	"github.com/sirdeggen/merkle/models"
	"github.com/sirdeggen/merkle/service"
)

type MerklePath []models.Hash

type MerklePathReader interface {
	Read(filename string, index uint64) (*MerklePath, error)
}

func Read(name string, index uint64) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make([]byte, 0, 8)

	// Read Unit64 from file
	UintBytes := make([]byte, 8)
	_, err = f.Read(UintBytes)
	if err != nil {
		return nil, err
	}
	// convert to uint64
	numOfTxs := binary.LittleEndian.Uint64(UintBytes)

	// is index within
	if index >= numOfTxs {
		return nil, fmt.Errorf("This block only has %d transactions, you tried to use index: %d which points to nil", numOfTxs, index)
	}

	// calulate how many power levels
	power := -1
	mask := uint64(1)
	for mask < uint64(numOfTxs) {
		mask = mask * 2
		power++
	}

	powerMask := uint64(mask)
	for powerMask < (math.MaxUint64 / 2) {
		powerMask = (powerMask * 2) | powerMask
	}

	// number of branches per level
	branches := make([]uint64, power+1)
	branches[0] = numOfTxs
	for x := 1; x <= power; x++ {
		branches[x] = uint64(math.Ceil(float64(branches[x-1]) / 2))
	}

	cumulativeBranchOffset := uint64(0)
	powerOffset := uint64(0)
	skip := uint64(0)
	for x := 0; x <= power; x++ {
		mask >>= 1
		branchOffset := branches[len(branches)-1-x]
		cumulativeBranchOffset += branchOffset
		if index&mask > 0 {
			powerOffset++
			// the tx is in the right branch
			// therefore we read the left path
			d := make([]byte, 32)
			_, err = f.Read(d)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
			// and add one to the offset
		} else {
			// the tx is in the left branch
			// therefore we read the right path by skipping forward one
			// unless it's the right most, in which case it's a duplicate of the left
			rightmost := cumulativeBranchOffset - 1
			if skip < rightmost {
				_, err = f.Seek(int64(32), 1)
				if err != nil {
					return nil, err
				}
			}
			d := make([]byte, 32)
			_, err = f.Read(d)
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}

		// calculate skip
		powerOffset <<= 1
		skip = powerOffset + cumulativeBranchOffset
		seekPosition := (32 * skip) + 8
		_, err = f.Seek(int64(seekPosition), 0)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func main() {
	block, err := service.GetBlockFromFile("data/midblock.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("There are ", len(block.Txids), " Transactions in this block.")
	branches, err := service.CalculateMerkleBranches(block)
	if err != nil {
		fmt.Println(err)
	}
	m := helpers.Reverse(block.MerkleRoot)
	fmt.Println("Merkle Root: ", hex.EncodeToString(m[:]))

	root := branches[len(branches)-1][0]
	cm := helpers.Reverse(root)
	fmt.Println("Calculated Merkle Root: ", hex.EncodeToString(cm[:]))

	block, err = service.GetBlockFromFile("data/midblock.json")
	if err != nil {
		fmt.Println(err)
	}

	// check all of the merkle paths
	for x := 0; x < len(block.Txids); x++ {
		txid := block.Txids[x]
		// read data/branches.bin
		data, err := Read("data/branches.bin", uint64(x))
		if err != nil {
			fmt.Println(err)
		}

		pathos := make([]models.Hash, 0, 0)
		for x := 0; x < 5; x++ {
			hash := [32]byte{}
			copy(hash[:], data[x*32:(x+1)*32])
			// prepend
			pathos = append([]models.Hash{hash}, pathos...)
		}
		path := &models.MerklePathBinary{
			Path:  pathos,
			Index: uint64(x),
		}
		valid := service.CheckMerklePathLeadsToRoot(&txid, path, &root)
		fmt.Println("Merkle Proof Valid: ", valid)
	}

}
