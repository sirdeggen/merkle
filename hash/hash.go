package hash

import (
	"encoding/hex"
	"fmt"

	"github.com/sirdeggen/merkle/helpers"
)

type Hash [32]byte

func (h *Hash) StringReverse() string {
	rev := helpers.Reverse(*h)
	return hex.EncodeToString(rev[:])
}

func FromStringReverse(s string) Hash {
	var h Hash
	bytes, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	copy(h[:], bytes)
	rev := helpers.Reverse(h)
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
