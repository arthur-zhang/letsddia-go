package lsm

import (
	"bytes"
	"encoding/binary"
	"io"
)

type KVPair struct {
	SeqId uint64
	Op    uint8
	Key   Key
	Value Value
}

func (kv *KVPair) size() uint64 {
	return 8 + 1 + 4 + uint64(len(kv.Key)) + 4 + uint64(len(kv.Value))
}
func (kv *KVPair) encode() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, kv.SeqId)
	_ = binary.Write(buf, binary.LittleEndian, kv.Op)
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(kv.Key)))
	_ = binary.Write(buf, binary.LittleEndian, []byte(kv.Key))
	_ = binary.Write(buf, binary.LittleEndian, uint32(len(kv.Value)))
	_ = binary.Write(buf, binary.LittleEndian, kv.Value)
	return buf.Bytes()
}
func decodeKvPair(r io.Reader) KVPair {
	var kv KVPair
	_ = binary.Read(r, binary.LittleEndian, &kv.SeqId)
	_ = binary.Read(r, binary.LittleEndian, &kv.Op)
	var keyLen uint32
	_ = binary.Read(r, binary.LittleEndian, &keyLen)

	bytes := make([]byte, keyLen)
	_ = binary.Read(r, binary.LittleEndian, bytes)
	kv.Key = Key(bytes)
	var valueLen uint32
	_ = binary.Read(r, binary.LittleEndian, &valueLen)
	kv.Value = make([]byte, valueLen)
	_ = binary.Read(r, binary.LittleEndian, &kv.Value)
	return kv
}

type DataBlock struct {
	KvList    []KVPair
	BlockSize uint64
}

const MaxDataBlockSize = uint64(128)

func NewDataBlock() DataBlock {
	return DataBlock{
		KvList:    make([]KVPair, 0),
		BlockSize: 0,
	}
}
func (b *DataBlock) kvCount() uint64 {
	return uint64(len(b.KvList))
}
func (b *DataBlock) hasSpace(kv *KVPair) bool {
	return b.BlockSize+kv.size() <= MaxDataBlockSize
}
func (b *DataBlock) Add(kv *KVPair) bool {
	if !b.hasSpace(kv) {
		return false
	}
	b.KvList = append(b.KvList, *kv)
	b.BlockSize += kv.size()
	return true
}
func (b *DataBlock) LastKv() *KVPair {
	if len(b.KvList) == 0 {
		return nil
	}
	return &b.KvList[len(b.KvList)-1]
}
func (b *DataBlock) Encode() []byte {
	buf := new(bytes.Buffer)
	for _, kv := range b.KvList {
		_ = binary.Write(buf, binary.LittleEndian, kv.encode())
	}
	return buf.Bytes()
}

func DecodeDataBlock(r io.Reader, blockSize uint64) DataBlock {
	var dataBlock DataBlock
	idx := uint64(0)
	for idx < blockSize {
		kv := decodeKvPair(r)
		dataBlock.KvList = append(dataBlock.KvList, kv)
		idx += kv.size()
	}
	dataBlock.BlockSize = blockSize
	return dataBlock
}
