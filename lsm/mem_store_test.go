package lsm

import (
	"commons-io/file_utils"
	"fmt"
	"math/rand"
	"testing"
)

func TestMemStore(t *testing.T) {
	baseDir := "/tmp/tiny-lsm"
	file_utils.Mkdirs(baseDir)
	file_utils.CleanDirectory(baseDir)
	diskStore := NewDiskStore(baseDir)
	config := Config{
		baseDir:          baseDir,
		maxMemStoreSize:  1024,
		blockSizeUpLimit: 128,
	}
	memStore := NewMemStore(&diskStore, &config)
	for i := 0; i < 1000; i++ {
		idx := rand.Int31n(100)
		key := Key(fmt.Sprintf("hello%09d", idx))
		value := Value(fmt.Sprintf("world%09d", idx))
		memStore.Add(uint64(i), 0, key, value)
	}
	memStore.flushIfNeed()
	fmt.Printf("%+v\n", memStore)

}
