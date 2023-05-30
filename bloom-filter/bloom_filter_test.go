package bloom_filter

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

type Key string

func (k Key) Bytes() []byte {
	return []byte(k)
}
func TestSingleElement(t *testing.T) {
	bf := NewBloomFilter[Key](10000, 0.10)
	item := Key("Hello")
	bf.Insert(item)
	if !bf.LookUp(item) {
		t.Errorf("Expected item to exist in the bloom filter")
	}
}

func TestMultipleElements(t *testing.T) {
	bf := NewBloomFilter[Key](10000, 0.10)
	items := []string{
		"Hello",
		"World",
		"Bloom",
		"Filter",
	}

	for _, item := range items {
		bf.Insert(Key(item))
		if !bf.LookUp(Key(item)) {
			t.Errorf("Expected item to exist in the bloom filter")
		}
	}
}
func TestNonExistentElement(t *testing.T) {
	bf := NewBloomFilter[Key](10000, 0.10)
	item := Key("Hello")
	nonExistentItem := Key("NonExistent")
	bf.Insert(item)
	if bf.LookUp(nonExistentItem) {
		t.Errorf("Expected item not to exist in the bloom filter")
	}
}

func TestBoundaryCondition(t *testing.T) {
	size := 10000
	bf := NewBloomFilter[Key](uint32(size), 0.10)

	for i := 0; i < size; i++ {
		bf.Insert(Key(strconv.Itoa(i)))
	}

	for i := 0; i < size; i++ {
		if !bf.LookUp(Key(strconv.Itoa(i))) {
			t.Errorf("Expected item to exist in the bloom filter")
		}
	}
}

func TestFalsePositiveRate(t *testing.T) {
	expectedFalsePositiveRate := 0.02
	size := 100_0000
	bf := NewBloomFilter[Key](uint32(size), expectedFalsePositiveRate)

	dataNum := 10000_0000_0000
	m := make(map[string]bool)
	for i := 0; i < size; i++ {
		num := rand.Int63n(int64(dataNum))
		s := fmt.Sprintf("%d", num)
		bf.Insert(Key(s))
		m[s] = true
	}

	falsePositiveCount := 0

	testCount := 1000000
	for i := 0; i < testCount; i++ {
		num := rand.Int63n(int64(dataNum))
		s := fmt.Sprintf("%d", num)
		_, found := m[s]

		if bf.LookUp(Key(s)) && !found {
			falsePositiveCount++
		}
	}

	falsePositiveRate := float64(falsePositiveCount) / float64(testCount)
	t.Logf("False positive rate: %f", falsePositiveRate)
	if math.Abs(falsePositiveRate-expectedFalsePositiveRate) > 0.01 {
		t.Errorf("False positive rate is too high")
	}
}

func TestDataSetSize(t *testing.T) {
	size := 100_0000
	bf := NewBloomFilter[Key](uint32(size), 0.02)
	t.Log(bf.m)
	bf = NewBloomFilter[Key](uint32(size*100), 0.02)
	t.Log(bf.m)
	bf = NewBloomFilter[Key](uint32(size/100), 0.02)
	t.Log(bf.m)
}
