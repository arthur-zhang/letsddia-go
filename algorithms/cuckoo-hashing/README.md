# Cuckoo Hashing implementation in Go

Cuckoo hashing is a scheme in computer programming for resolving hash collisions of values of hash functions in a table, with worst-case constant lookup time. The name derives from the behavior of some species of cuckoo, where the cuckoo chick pushes the other eggs or young out of the nest when it hatches in a variation of the behavior referred to as brood parasitism; analogously, inserting a new key into a cuckoo hashing table may push an older key to a different location in the table.

### Demo with two hash function

```go
func hash0(key int) int {
	return key % 10
}
func hash1(key int) int {
	return key / 10
}
```

| n  | hash0%SIZE | hash1%SIZE |
|----|------------|------------|
| 10 | 0          | 1          |
| 20 | 0          | 2          |
| 30 | 0          | 3          |
| 40 | 0          | 4          |
| 50 | 0          | 0          |

### After insert

| data0 | data0 | data1 |
|-------|-------|-------|
| 0     | 50    | 40    |
| 1     | -1    | 10    |
| 2     | -1    | 20    |
| 3     | -1    | 30    |


### Usage

```go
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
```

    