package simhash

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"strings"
)

func hash_(input string) uint64 {
	h := murmur3.New64()
	_, _ = h.Write([]byte(input))
	hash := h.Sum64()
	return hash
}

func hammingDistance(x, y uint64) int {
	var dist int
	val := x ^ y
	fmt.Printf("%064b\n", val)
	for val != 0 {
		dist++
		// clear the least significant bit set
		val &= val - 1
	}
	return dist
}
func SimHash(input string) uint64 {
	words := strings.Split(input, " ")
	bits := make([]int, 64)
	for _, word := range words {
		v := hash_(word)
		for i := 64; i >= 1; i-- {
			if (v>>(64-i))&0x1 == 1 {
				bits[i-1]++
			} else {
				bits[i-1]--
			}
		}
	}
	var res uint64
	for i := 64; i >= 1; i-- {
		if bits[i-1] > 0 {
			res |= 1 << (64 - i)
		}
	}
	return res
}
