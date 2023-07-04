package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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
	fmt.Println("numOfTxs: ", numOfTxs)

	// is index within
	if index >= numOfTxs {
		return nil, fmt.Errorf("This block only has %d transactions, you tried to use index: %d which points to nil", numOfTxs, index)
	}

	// flip each bit of index
	inverseIndex := ^index

	// calulate how many power levels
	power := -1
	mask := uint64(1)
	for mask < uint64(numOfTxs) {
		mask = mask * 2
		power++
	}

	skip := 2
	for x := 0; x <= power; x++ {
		mask >>= 1
		fmt.Printf("%64b\n", mask)
		fmt.Printf("%64b\n", inverseIndex)
		if inverseIndex&mask > 0 {
			fmt.Println("r")
			d := make([]byte, 32)
			_, err = f.Seek(int64(32), 1)
			if err != nil {
				return nil, err
			}
			_, err = f.Read(d)
			var t [32]byte
			copy(t[:], d)
			hash := helpers.Reverse(t)
			fmt.Println("hash: ", hex.EncodeToString(hash[:]))
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		} else {
			fmt.Println("l")
			skip += 1
			d := make([]byte, 32)
			_, err = f.Read(d)
			var t [32]byte
			copy(t[:], d)
			hash := helpers.Reverse(t)
			fmt.Println("hash: ", hex.EncodeToString(hash[:]))
			if err != nil {
				return nil, err
			}
			data = append(data, d...)
		}
		// calculate skip
		skip = skip * 2
		fmt.Println("skip: ", skip)
		_, err = f.Seek(int64(skip), 1)
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

	cm := helpers.Reverse(branches[len(branches)-1][0])
	fmt.Println("Calculated Merkle Root: ", hex.EncodeToString(cm[:]))

	// read data/branches.bin
	data, err := Read("data/branches.bin", 10)
	if err != nil {
		fmt.Println(err)
	}
	// print each line as hex string
	for x := 1; x <= len(data)/32; x++ {
		//data[(x-1)*32:x*32]
		// reverse the 32 bytes
		var chunk [32]byte
		copy(chunk[:], data[(x-1)*32:x*32])
		d := helpers.Reverse(chunk)
		fmt.Println("Merkle Path: ", hex.EncodeToString(d[:]))
	}

	block, err = service.GetBlockFromFile("data/midblock.json")
	if err != nil {
		fmt.Println(err)
	}

	pathos := make([]models.Hash, 0, 0)
	for x := 0; x < 5; x++ {
		fmt.Println("x: ", x, data[x*32:(x+1)*32])
		hash := [32]byte{}
		copy(hash[:], data[x*32:(x+1)*32])
		pathos = append(pathos, hash)
	}

	path := &models.MerklePathBinary{
		Path:  pathos,
		Index: 10,
	}

	txid := block.Txids[10]
	rev := helpers.Reverse(txid)
	fmt.Println("Txid: ", hex.EncodeToString(rev[:]))

	valid := service.CheckMerklePathLeadsToRoot(&txid, path, &branches[len(branches)-1][0])
	fmt.Println("Merkle Proof Valid: ", valid)

}
