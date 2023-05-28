package lsm

import (
	"testing"
)

func TestDiskFileIterator(t *testing.T) {
	diskFile := openDiskFile("/tmp/tiny-lsm/0000000001.dat")
	iter := diskFile.Iter()
	iter.SeekToFirst()
	for iter.Valid() {
		key := iter.Key()
		value := iter.Value()
		t.Logf("%s, %s\n", *key, *value)
		iter.Next()
	}
}

func TestSeek(t *testing.T) {
	diskFile := openDiskFile("/tmp/tiny-lsm/0000000001.dat")
	iter := diskFile.Iter()
	key := Key("hello000000078")
	for iter.Seek(&key); iter.Valid(); iter.Next() {
		t.Logf("%s, %s\n", *iter.Key(), *iter.Value())
	}
}
