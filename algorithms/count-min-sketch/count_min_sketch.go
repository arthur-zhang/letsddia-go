package count_min_sketch

import (
	"github.com/spaolacci/murmur3"
	"hash"
)

type CountMinSketch struct {
	d      int // d rows
	w      int // w columns
	table  [][]int
	hashes []hash.Hash32
}

func New(d int, w int) *CountMinSketch {
	table := make([][]int, d)
	for i := range table {
		table[i] = make([]int, w)
	}
	hashes := make([]hash.Hash32, d)
	for i := 0; i < d; i++ {
		hashes[i] = murmur3.New32WithSeed(uint32(i))
	}
	return &CountMinSketch{d: d, w: w, table: table, hashes: hashes}
}
func (cms *CountMinSketch) hash(item string, i int) uint32 {
	h := cms.hashes[i]
	h.Reset()
	_, _ = h.Write([]byte(item))
	j := h.Sum32() % uint32(cms.w)
	return j
}
func (cms *CountMinSketch) Update(item string, count int) {
	for i := 0; i < cms.d; i++ {
		j := cms.hash(item, i)
		cms.table[i][j] += count
	}
}
func (cms *CountMinSketch) Estimate(item string) int {
	min := 1 << 31 // int32 max
	for i := 0; i < cms.d; i++ {
		j := cms.hash(item, i)
		if cms.table[i][j] < min {
			min = cms.table[i][j]
		}
	}
	return min
}
