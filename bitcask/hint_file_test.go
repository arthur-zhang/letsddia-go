package bitcask

import (
	"commons-io/file_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPut(t *testing.T) {
	file := OpenHintFile("/tmp/hint_file", file_utils.NewOpenOptions().Read(true).Write(true).Create(true))
	file.put("hello", KeyDirEntry{
		fileId:   1,
		valueSz:  10,
		valuePos: 0,
		tstamp:   1234,
	})

	iter := file.NewIterator()
	for item := iter.Next(); item != nil; item = iter.Next() {
		println(item.key)
	}
}

func TestPutAndRead(t *testing.T) {
	datFilePath := "/tmp/bitcask/0000000148.dat"
	datFile := newDatFileBuilder().path(datFilePath).openOptions(file_utils.NewOpenOptions().Read(true)).build()
	datIter := datFile.NewIterator()
	hintFilePath := "/tmp/bitcask/0000000148.hint"
	hintFile := OpenHintFile(hintFilePath, file_utils.NewOpenOptions().Read(true))
	hintIter := hintFile.NewIterator()
	for {
		datItem := datIter.Next()
		hintItem := hintIter.Next()

		if datItem == nil && hintItem == nil {
			break
		}
		assert.Equal(t, datItem.block.key, hintItem.key)
	}
}
