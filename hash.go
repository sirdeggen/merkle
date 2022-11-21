package main

import (
	"encoding/hex"
	"fmt"
)

type Hash [32]byte

func hexToBytes(h string) ([]byte, error) {
	bytes, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func reverse(b [32]byte) [32]byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func (h *Hash) toStringReverse() string {
	rev := reverse(*h)
	return hex.EncodeToString(rev[:])
}

func fromStringReverse(s string) Hash {
	var h Hash
	bytes, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	copy(h[:], bytes)
	rev := reverse(h)
	return rev
}

func (h *Hash) toString() string {
	return hex.EncodeToString(h[:])
}

func fromString(s string) Hash {
	var h Hash
	bytes, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	copy(h[:], bytes)
	return h
}
