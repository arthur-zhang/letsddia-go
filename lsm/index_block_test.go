package lsm

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndexBlock_encode(t *testing.T) {
	indexBlock := NewIndexBlock()
	for i := 0; i < 10; i++ {
		meta := NewBlockMeta(KVPair{1, 2, "key", []byte("value")}, 1, 2)
		indexBlock.appendBlockIndex(meta)
	}
	encoded := indexBlock.encode()
	decoded := decodeIndexBlock(bytes.NewBuffer(encoded), indexBlock.totalSize)
	assert.Equal(t, indexBlock, decoded)
}
