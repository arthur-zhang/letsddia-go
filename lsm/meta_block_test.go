package lsm

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetaBlock_encode(t *testing.T) {
	metaBlock := NewMetaBlock(100, 10, 100, 100)
	encoded := metaBlock.encode()
	decoded := decodeMetaBlock(bytes.NewBuffer(encoded))
	assert.Equal(t, metaBlock, decoded)
}
