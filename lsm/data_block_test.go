package lsm

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataBlock_Encode(t *testing.T) {
	dataBlock := NewDataBlock()
	for i := 0; i < 10; i++ {
		kv := KVPair{
			Key:   Key(fmt.Sprintf("hello%09d", i)),
			Value: []byte(fmt.Sprintf("world%09d", i)),
		}
		dataBlock.Add(&kv)
	}
	encodeData := dataBlock.Encode()

	decoded := DecodeDataBlock(bytes.NewBuffer(encodeData), dataBlock.BlockSize)
	assert.Equal(t, dataBlock, decoded)
}

func TestKvEncodeDecode(t *testing.T) {
	kv := KVPair{
		SeqId: 1,
		Op:    1,
		Key:   "hello",
		Value: []byte("world"),
	}
	encoded := kv.encode()
	decodedKv := decodeKvPair(bytes.NewBuffer(encoded))
	assert.Equal(t, kv, decodedKv)
}
