package tree

import (
	"math"
	"testing"
)

func BenchmarkTestReadSmall(b *testing.B) {
	mts := NewMerkleTreeService("../data")
	for i := 0; i < b.N; i++ {
		fl := float64(i)
		x := uint64(math.Mod(fl, 1))
		_, err := mts.Read("0a3b8cb97063d49e1a1504f10c5c6e648ec8fc436f1ab0ee68dd457a305f0a8b", x)
		if err != nil {
			break
		}
	}
}

func BenchmarkTestReadMedium(b *testing.B) {
	mts := NewMerkleTreeService("../data")
	for i := 0; i < b.N; i++ {
		fl := float64(i)
		x := uint64(math.Mod(fl, 17))
		_, err := mts.Read("d66e56fb408763e36e8622eb56a8a1072ccc606476fe9e0765cca0dff95949b1", x)
		if err != nil {
			break
		}
	}
}

func BenchmarkTestReadLarge(b *testing.B) {
	mts := NewMerkleTreeService("../data")
	for i := 0; i < b.N; i++ {
		fl := float64(i)
		x := uint64(math.Mod(fl, 2713))
		_, err := mts.Read("a623039e2030dfafd02af3948f8f8483aadbb7296e205ff95f78a52269be97f5", x)
		if err != nil {
			break
		}
	}
}

func BenchmarkTestReadHuge(b *testing.B) {
	mts := NewMerkleTreeService("../data")
	for i := 0; i < b.N; i++ {
		fl := float64(i)
		x := uint64(math.Mod(fl, 27731))
		_, err := mts.Read("4cb1de606d95d4b359a20f0e26e4227383b12ca759a9af8e6423f89cb33f850a", x)
		if err != nil {
			break
		}
	}
}
