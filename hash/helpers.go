package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Hash [32]byte

func Reverse(b [32]byte) [32]byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func HexToBytes(h string) ([]byte, error) {
	bytes, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (h *Hash) StringReverse() string {
	rev := Reverse(*h)
	return hex.EncodeToString(rev[:])
}

func FromStringReverse(s string) Hash {
	var h Hash
	bytes, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	copy(h[:], bytes)
	rev := Reverse(h)
	return rev
}

func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}

func FromString(s string) Hash {
	var h Hash
	bytes, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	copy(h[:], bytes)
	return h
}

func Sha256Sha256(digest []byte) [32]byte {
	one := sha256.Sum256(digest)
	return sha256.Sum256(one[:])
}
