package hyperloglog

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func TestTrailingZeros(t *testing.T) {
	// Test a few simple cases
	testCases := []struct {
		input uint64
		want  int
	}{
		{input: 0, want: 64},
		{input: 1, want: 0},
		{input: 2, want: 1},
		{input: 3, want: 0},
		{input: 4, want: 2},
		{input: 8, want: 3},
	}

	for _, tc := range testCases {
		got := trailingZeros(tc.input)
		if got != tc.want {
			t.Errorf("trailingZeros(%d) = %d; want %d", tc.input, got, tc.want)
		}
	}
}
func TestHyperLogLog(t *testing.T) {
	hll := New(14)

	println(hll.m)
	// Add 10000 distinct items
	numItems := 10000
	for i := 0; i < numItems; i++ {
		hll.Add(fmt.Sprintf("%d", i))
	}

	// The count should be approximately equal to the number of items added
	// Allow 1% error
	count := hll.Cardinality()
	println(uint64(count))
	if count < float64(numItems)*0.99 || count > float64(numItems)*1.01 {
		t.Errorf("Expected count to be approximately %d, got %f", numItems, count)
	}

	// Adding the same item multiple times shouldn't change the count
	hll.Add("0")
	countAfterAddingSameItem := hll.Cardinality()
	if count != countAfterAddingSameItem {
		t.Errorf("Expected count to remain the same after adding the same item, but got %f", countAfterAddingSameItem)
	}
}
func TestHyperLogLog_Cardinality(t *testing.T) {
	hll := New(16)
	hll.AddInt(0)
	hll.AddInt(1)
	hll.AddInt(2)
	hll.AddInt(3)
	hll.AddInt(16)
	hll.AddInt(17)
	hll.AddInt(18)
	hll.AddInt(19)
	hll.AddInt(19)
	card := uint64(hll.Cardinality())
	assert.Equal(t, 8, card)
}

func TestHighCardinality(t *testing.T) {
	hll := New(14)
	size := 100_0000
	set := make(map[string]struct{})
	for i := 0; i < size; i++ {
		n := rand.Uint64()
		str := fmt.Sprintf("%x", n)
		hll.Add(str)
		set[str] = struct{}{}
	}
	card := hll.Cardinality()

	setSize := len(set)
	fmt.Printf("card %f, setSize: %d\n", card, setSize)
	fmt.Printf("diff: %f\n", math.Abs(card-float64(setSize)))
	err := math.Abs(card-float64(setSize)) / float64(setSize)
	fmt.Printf("err: %f\n", err)
	assert.True(t, err < 0.01)

}
