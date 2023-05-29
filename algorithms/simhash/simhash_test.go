package simhash

import (
	"fmt"
	"testing"
)

func TestSimHash(t *testing.T) {
	testCases := []struct {
		text1     string
		text2     string
		threshold int
	}{
		{"Hello, World!", "Hello, World!", 0},
		{"Hello, World!", "Hello, world!", 1},
		{"Hello, World!", "Goodbye, World!", 2},
		{"Hello, World!", "Completely different text", 6},
		{"The quick brown fox jumps over the lazy dog.", "The quick brown fox jumps over the lazy cat.", 4},
	}
	for _, tc := range testCases {
		hash1 := SimHash(tc.text1)
		hash2 := SimHash(tc.text2)
		dist := hammingDistance(hash1, hash2)
		fmt.Printf("hamming distance between %s and %s is %d\n", tc.text1, tc.text2, dist)
		if dist > tc.threshold {
			t.Errorf("hamming distance between %s and %s is %d, expected %d", tc.text1, tc.text2, dist, tc.threshold)
		}
	}
}
