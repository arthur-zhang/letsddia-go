package lsm

import (
	"commons-io/file_utils"
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	println("......")
}

func TestNew(t *testing.T) {
	config := Config{
		baseDir:          "/tmp/tiny-lsm",
		maxMemStoreSize:  1024,
		blockSizeUpLimit: 128,
	}
	_ = file_utils.CleanDirectory(config.baseDir)
	lsm := New(config)
	for i := 0; i < 1000; i++ {
		//idx := rand.Int31n(100)
		idx := i % 100
		key := Key(fmt.Sprintf("hello%09d", idx))
		value := Value(fmt.Sprintf("world%09d", idx))
		lsm.Put(key, value)
	}
	key := Key("hello000000078")

	v, found := lsm.Get(&key)
	if found {
		t.Logf("%s\n", string(v))
	} else {
		t.Logf("not found\n")
	}
}
func TestScan(t *testing.T) {

	config := Config{
		baseDir:          "/tmp/tiny-lsm",
		maxMemStoreSize:  1024,
		blockSizeUpLimit: 128,
	}
	_ = file_utils.CleanDirectory(config.baseDir)
	lsm := New(config)
	for i := 0; i < 1000; i++ {
		idx := i % 100
		key := Key(fmt.Sprintf("hello%09d", idx))
		value := Value(fmt.Sprintf("world%09d", idx))
		lsm.Put(key, value)
	}
	startKey := Key("hello000000078")
	//startKey := Key("hello000000078")
	endKey := Key("hello000000081")

	iter := lsm.Scan(&startKey, &endKey)
	for iter.Valid() {
		t.Logf("%s, %s\n", *iter.Key(), *iter.Value())
		iter.Next()
	}

}
