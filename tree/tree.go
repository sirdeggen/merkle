package tree

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"os"

	"github.com/sirdeggen/merkle/hash"
)

type MerkleTree [][]hash.Hash
type MerklePath []hash.Hash
type MerklePathJson struct {
	Item  string   `json:"item"`
	Index string   `json:"index"`
	Path  []string `json:"path"`
}

type MerkleTreeReaderWriter interface {
	Read(string, uint64) (*MerklePathJson, error)
	Write(string, MerkleTree) error
}

// merkleTreeService is both a MerkleTreeReader and MerkleTreeWriter
type merkleTreeService struct {
	Directory string
}

func NewMerkleTreeService(dir string) *merkleTreeService {
	return &merkleTreeService{
		Directory: dir,
	}
}

func (mpw *merkleTreeService) Write(branches MerkleTree) error {
	root := branches[0][0].StringReverse()
	l := branches[len(branches)-1]
	numOfTxs := make([]byte, 8)
	binary.LittleEndian.PutUint64(numOfTxs, uint64(len(l)))
	fileBytes := numOfTxs
	for x := len(branches) - 2; x >= 0; x-- {
		for y := 0; y < len(branches[x]); y++ {
			fileBytes = append(fileBytes, branches[x][y][:]...)
		}
	}
	err := ioutil.WriteFile(fmt.Sprint(mpw.Directory, '/', root), fileBytes, fs.FileMode(0644))
	return err
}

func (mpw *merkleTreeService) Read(root string, index uint64) (*MerklePathJson, error) {
	filename := fmt.Sprint(mpw.Directory, '/', root)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var path MerklePath

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
	var item hash.Hash
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
			var h hash.Hash
			copy(h[:], d)
			path = append(path, h)
			if x == power {
				// we are at the leaf level
				// read the item
				d := make([]byte, 32)
				_, err = f.Read(d)
				if err != nil {
					return nil, err
				}
				copy(item[:], d)
			}
		} else {
			// the tx is in the left branch
			// therefore we read the right path by skipping forward one
			// unless it's the right most, in which case it's a duplicate of the left so we don't skip
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
			var h hash.Hash
			copy(h[:], d)
			path = append(path, h)
			if x == power {
				// we are at the leaf level
				// read the item which is two back
				backTwo := -2 * int64(32)
				_, err = f.Seek(backTwo, 1)
				if err != nil {
					return nil, err
				}
				d := make([]byte, 32)
				_, err = f.Read(d)
				if err != nil {
					return nil, err
				}
				copy(item[:], d)
			}
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
	var merklePath MerklePathJson
	for _, h := range path {
		merklePath.Path = append(merklePath.Path, h.StringReverse())
	}
	merklePath.Item = item.StringReverse()
	return &merklePath, nil
}
