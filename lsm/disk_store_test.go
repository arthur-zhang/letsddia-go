package lsm

import (
	"testing"
)

func TestDiskStoreIter(t *testing.T) {

	store := NewDiskStore("/tmp/tiny-lsm")
	iter := store.Iter()
	iter.SeekToFirst()
	for iter.Valid() {
		key := iter.Key()
		value := iter.Value()
		t.Logf("%s, %s\n", *key, *value)
		iter.Next()
	}
	println(">>>>")
}
