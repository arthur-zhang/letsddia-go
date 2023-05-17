package cuckoo_hashing

import "testing"

func hash0(key int) int {
	return key % 10
}
func hash1(key int) int {
	return key / 10
}

func TestCuckooHashing(t *testing.T) {

	cht := New(4, hash0, hash1)

	// insert key
	t.Run("Insert", func(t *testing.T) {
		keys := []int{10, 20, 30, 40, 50}
		for _, key := range keys {
			cht.Insert(key)
		}

		// ensure all keys are present
		for _, key := range keys {
			if ok := cht.Lookup(key); !ok {
				t.Errorf("Expected key %d to be present, but it was not", key)
			}
		}
	})

	// lookup key
	t.Run("Lookup", func(t *testing.T) {
		if ok := cht.Lookup(10); !ok {
			t.Error("Expected key 10 to be present, but it was not")
		}

		if ok := cht.Lookup(100); ok {
			t.Error("Expected key 100 to be absent, but it was found")
		}
	})

	// delete key
	t.Run("Erase", func(t *testing.T) {
		cht.Erase(30)
		if ok := cht.Lookup(30); ok {
			t.Error("Expected key 30 to be absent after erasing, but it was found")
		}
	})

	// delete nonexistent key
	t.Run("Erase nonexistent key", func(t *testing.T) {
		cht.Erase(100)
		if ok := cht.Lookup(100); ok {
			t.Error("Expected key 100 to be absent, but it was found")
		}
	})
}
