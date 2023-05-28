package lsm

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockMeta_encode(t *testing.T) {
	meta := NewBlockMeta(KVPair{1, 2, "key", []byte("value")}, 1, 2)
	encoded := meta.encode()
	decoded := decodeBlockMeta(bytes.NewBuffer(encoded))
	assert.Equal(t, meta, decoded)
}
