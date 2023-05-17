package cuckoo_hashing

type CuckooHashingTable struct {
	data0 []int
	data1 []int
	hash0 HashFunc
	hash1 HashFunc
	size  int
}

func New(n int, hash0, hash1 HashFunc) *CuckooHashingTable {
	data0 := make([]int, n)
	data1 := make([]int, n)
	for i := 0; i < n; i++ {
		data0[i] = -1
		data1[i] = -1
	}
	return &CuckooHashingTable{
		data0: data0,
		data1: data1,
		hash0: hash0,
		hash1: hash1,
		size:  n,
	}
}

type HashFunc func(key int) int

func (cht *CuckooHashingTable) Lookup(key int) bool {

	h1 := cht.hash0(key) % cht.size
	if cht.data0[h1] == key {
		return true
	}
	h2 := cht.hash1(key) % cht.size
	if cht.data1[h2] == key {
		return true
	}
	return false
}

func (cht *CuckooHashingTable) Erase(key int) bool {

	h1 := cht.hash0(key) % cht.size
	if cht.data0[h1] == key {
		cht.data0[h1] = -1
		return true
	}
	h2 := cht.hash1(key) % cht.size
	if cht.data1[h2] == key {
		cht.data1[h2] = -1
		return true
	}
	return false
}
func (cht *CuckooHashingTable) Insert(key int) {
	cht.insertImpl(key, 0, 0)
}
func (cht *CuckooHashingTable) insertImpl(key int, loopCount int, tableIdx int) bool {
	if loopCount >= cht.size {
		println("Cycle detected, while inserting ", key, ", rehashing required")
		return false
	}

	if tableIdx == 0 {
		h1 := cht.hash0(key) % cht.size
		if cht.data0[h1] == -1 {
			cht.data0[h1] = key
			println("Inserted key ", key, " at table1")
			return true
		}
		// store the old key
		old := cht.data0[h1]
		// insert the new key to data0
		cht.data0[h1] = key
		println("Inserted key ", key, " at table1")
		// recursively insert the old key
		return cht.insertImpl(old, loopCount+1, 1)
	} else {
		h2 := cht.hash1(key) % cht.size
		if cht.data1[h2] == -1 {
			cht.data1[h2] = key
			println("Inserted key ", key, " at table2")
			return true
		}
		// store the old key
		old := cht.data1[h2]
		println("Inserted key ", key, " at table2")
		// insert the new key to data0
		cht.data1[h2] = key
		return cht.insertImpl(old, loopCount+1, 0)
	}
}
