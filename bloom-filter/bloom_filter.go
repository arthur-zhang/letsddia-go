package bloom_filter

import (
	"bitset"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash"
	"hash/crc64"
	"hash/fnv"
	"math"
)

type Byteable interface {
	Bytes() []byte
}

// BloomFilter Bloom filters have two main components:
// 1. A bit array of m bits, initially all set to 0.
// 2. k different hash functions h1(x), h2(x), ..., hk(x).
// To insert an element x into the filter, we compute the k hash values h1(x), h2(x), ..., hk(x) and set the bits at these positions to 1.
// To query whether an element y is in the filter, we compute the k hash values h1(y), h2(y), ..., hk(y) and check whether the bits at these positions are all 1.
type BloomFilter[T Byteable] struct {
	n    uint32        // number of items in the filter
	f    float64       // false positive rate
	m    uint32        // bloom filter size(bits) m = -n*ln(f)/(ln(2))^2
	k    int           // number of hash functions  k = m/n*ln(2)
	bits bitset.BitSet // A bit array of m bits, initially all set to 0.
	// k different hash functions h1(x), h2(x), ..., hk(x).
	hashFuncs []HashFunc[T]
}

type HashFunc[T Byteable] func(item T) uint64

func HashAdapter[T Byteable](h hash.Hash) HashFunc[T] {
	return func(item T) uint64 {
		h.Reset()
		h.Write(item.Bytes())
		bytes := h.Sum(nil)
		return binary.BigEndian.Uint64(bytes)
	}
}
func (bf *BloomFilter[T]) getHashFuncs(k int) []HashFunc[T] {
	funcs := DefaultHashFuncs[T]()
	if k >= len(funcs) {
		return funcs
	}
	return funcs[:k]
}
func DefaultHashFuncs[T Byteable]() []HashFunc[T] {
	return []HashFunc[T]{
		HashAdapter[T](murmur3.New64()),
		HashAdapter[T](crc64.New(crc64.MakeTable(crc64.ECMA))),
		HashAdapter[T](fnv.New64()),
		HashAdapter[T](fnv.New128()),
		HashAdapter[T](md5.New()),
		HashAdapter[T](sha1.New()),
	}
}

// calcM Compute bitmap size for n items and false positive rate f
func calcM(n uint32, f float64) uint32 {
	if f <= 0 || f >= 1 {
		panic("f must be in [0, 1)")
	}
	// m = -n*ln(f)/(ln(2))^2
	return uint32(math.Ceil(-float64(n) * math.Log(f) / math.Pow(math.Log(2), 2)))
}

func calcK(m, n uint32) int {
	// k = m/n*ln(2)
	k := int(math.Ceil(float64(m) / float64(n) * math.Log(2)))
	if k < 1 {
		k = 1
	}
	return k
}

// calcN round up n to the nearest multiple of 8
func calcN(n uint32) uint32 {
	n = (n + 7) &^ 7
	return n >> 3
}

func NewBloomFilter[T Byteable](n uint32, f float64) *BloomFilter[T] {
	m := calcM(n, f)
	k := calcK(m, n)
	bits := bitset.New(m)
	m = bits.Len()

	println("n:", n, "f:", f, "m:", m, "k:", k)
	return &BloomFilter[T]{
		n:         n,
		f:         f,
		k:         k,
		m:         m,
		bits:      bits,
		hashFuncs: DefaultHashFuncs[T](),
	}
}

func (bf *BloomFilter[T]) Insert(item T) {
	for _, hashFunc := range bf.hashFuncs {
		h := hashFunc(item) % uint64(bf.m)
		bf.bits.Set(uint32(h))
	}
}
func (bf *BloomFilter[T]) LookUp(item T) bool {
	for _, hashFunc := range bf.hashFuncs {
		h := hashFunc(item) % uint64(bf.m)
		if !bf.bits.IsSet(uint32(h)) {
			return false
		}
	}
	return true
}

func (bf *BloomFilter[T]) DebugPrint() {
	for i, bit := range bf.bits {
		if (i % 16) == 0 {
			fmt.Println()
		}
		fmt.Printf("%08b ", bit)
	}
}
