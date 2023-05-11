package skiplist

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
)

func TestRandomLevel(t *testing.T) {
	s := make(map[int]int)
	for i := 0; i < 10000; i++ {
		level := randomLevel()
		s[level]++
	}
	for i := range s {
		t.Log(i, s[i])
	}
}
func TestEmpty(t *testing.T) {
	list := NewSkipList[int, int]()
	key := 10
	assert.False(t, list.Contains(&key))

	iter := list.Iterator()
	assert.False(t, iter.Valid())

	iter.SeekToFirst()
	assert.False(t, iter.Valid())

	key = 100
	iter.Seek(&key)
	assert.False(t, iter.Valid())

	iter.SeekToLast()
	assert.False(t, iter.Valid())
}

func setToSortedArray(set map[int]bool) []int {
	keys := make([]int, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
func TestInsertAndLookup(t *testing.T) {
	const N = 2000
	const R = 5000
	list := NewSkipList[int, int]()
	keys := make(map[int]bool)
	var zeroV int
	for i := 0; i < N; i++ {
		key := rand.Intn(R)
		// true if the key exists in the map
		_, exists := keys[key]
		if !exists {
			keys[key] = true
			list.Insert(&key, &zeroV)
		}
	}
	for i := 0; i < R; i++ {
		_, exists := keys[i]
		if list.Contains(&i) {
			assert.True(t, exists)
		} else {
			assert.False(t, exists)
		}
	}

	sortedKeys := setToSortedArray(keys)
	// Simple iterator tests
	{
		iter := list.Iterator()
		assert.False(t, iter.Valid())

		key := 0
		iter.Seek(&key)
		assert.True(t, iter.Valid())
		assert.Equal(t, sortedKeys[0], *iter.Key())

		iter.SeekToFirst()
		assert.True(t, iter.Valid())
		assert.Equal(t, sortedKeys[0], *iter.Key())

		iter.SeekToLast()
		assert.True(t, iter.Valid())
		assert.Equal(t, sortedKeys[len(sortedKeys)-1], *iter.Key())

	}
	// Forward iteration test
	for i := 0; i < R; i++ {
		iter := list.Iterator()
		iter.Seek(&i)

		index := sort.Search(len(sortedKeys), func(j int) bool { return sortedKeys[j] >= i })

		for j := 0; j < 3; j++ {
			if index == len(sortedKeys) {
				assert.False(t, iter.Valid())
				break
			} else {
				assert.True(t, iter.Valid())
				assert.Equal(t, sortedKeys[index], *iter.Key())
				iter.Next()
				index++
			}
		}
	}
}
