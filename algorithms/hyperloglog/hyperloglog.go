package hyperloglog

import (
	"bytes"
	"encoding/binary"
	"github.com/spaolacci/murmur3"
	"math"
)

func alphaM(b, m int) float64 {
	switch b {
	case 4:
		return 0.673 * float64(m*m)
	case 5:
		return 0.697 * float64(m*m)
	case 6:
		return 0.709 * float64(m*m)
	default:
		return 0.7213 / (1 + 1.079/float64(m))
	}
}
func firstBit(h uint64, b int) uint64 {
	return h >> (64 - b)
}
func trailingZeros(h uint64) int {
	count := 0
	for i := 0; i < 64; i++ {
		if h&1 == 1 {
			break
		}
		count++
		h >>= 1
	}
	return count
}

type HyperLogLog struct {
	b         int   // number of bits used to index into a bucket
	m         int   // number of buckets
	registers []int // array of m entries, storing max trailing zeros per bucket
	alphaM    float64
}

func New(b int) *HyperLogLog {
	m := 1 << b
	return &HyperLogLog{
		b:         b,
		m:         m,
		registers: make([]int, m),
		alphaM:    alphaM(b, m),
	}
}
func hash(x []byte) uint64 {
	h := murmur3.New64()
	_, _ = h.Write(x)
	return h.Sum64()
}
func hashInt(x int32) uint64 {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, x)
	return hash(buf.Bytes())
}
func (hll *HyperLogLog) Add(x string) {
	h := hash([]byte(x))
	p := trailingZeros(h) + 1
	bucket := firstBit(h, hll.b)
	if p > hll.registers[bucket] {
		hll.registers[bucket] = p
	}
}

func (hll *HyperLogLog) AddInt(x int32) {

	h := hashInt(x)
	p := trailingZeros(h) + 1
	bucket := firstBit(h, hll.b)
	if p > hll.registers[bucket] {
		hll.registers[bucket] = p
	}
}

func (hll *HyperLogLog) Cardinality() float64 {
	sum := 0.0
	zeros := 0
	for _, it := range hll.registers {
		sum += 1.0 / float64(int(1)<<it)
		if it == 0 {
			zeros++
		}
	}
	harmonicAvg := float64(hll.m) / sum
	E := hll.alphaM * float64(hll.m) * harmonicAvg // Raw estimate

	// Computes corrected estimate
	if E <= 2.5*float64(hll.m) { // Small range correction
		if zeros != 0 {
			return float64(hll.m) * math.Log(float64(hll.m)/float64(zeros))
		} else {
			return E
		}
	} else if E <= math.Pow(2, 32)/float64(30) { // Intermediate range, no correction
		return E
	} else { // E > 2^32/30, Large range correction
		return -math.Pow(2, 32) * math.Log(1-E/math.Pow(2, 32))
	}
}
