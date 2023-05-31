package bitset

import "math/bits"

const BITS = 64

type BitSet struct {
	data []uint64
	len  uint64
}

func New(bits uint64) BitSet {
	blocks, rem := divRem(bits)
	if rem > 0 {
		blocks++
	}
	return BitSet{make([]uint64, blocks), bits}
}
func divRem(x uint64) (uint64, uint64) {
	return x / BITS, x % (BITS - 1)
}
func (bs *BitSet) Len() uint64 {
	return bs.len
}

func (bs *BitSet) Count() int {
	count := 0
	for _, n := range bs.data {
		count += bits.OnesCount64(n)
	}
	return count
}

func (bs BitSet) Clear(bit uint64) bool {
	if bit >= bs.Len() {
		panic("out of range")
	}
	block, i := divRem(bit)
	word := bs.data[block]
	old := word&(1<<i) != 0
	bs.data[block] &= ^(1 << i)
	return old
}

// Set Enable bit and return old value
func (bs BitSet) Set(bit uint64) bool {
	if bit >= bs.Len() {
		panic("out of range")
	}
	block, i := divRem(bit)
	word := bs.data[block]
	old := word&(1<<i) != 0
	bs.data[block] |= 1 << i
	return old
}
func (bs BitSet) IsSet(n uint64) bool {
	if n >= bs.Len() {
		panic("out of range")
	}
	block, i := divRem(n)
	return bs.data[block]&(1<<i) != 0
}
