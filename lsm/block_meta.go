package lsm

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BlockMeta struct {
	lastKv    KVPair
	offset    uint64
	blockSize uint64
}

func NewBlockMeta(lastKv KVPair, offset, blockSize uint64) BlockMeta {
	return BlockMeta{
		lastKv:    lastKv,
		offset:    offset,
		blockSize: blockSize,
	}
}
func (m *BlockMeta) size() uint64 {
	return 8 + 8 + m.lastKv.size()
}

func (m *BlockMeta) encode() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, m.lastKv.encode())
	_ = binary.Write(buf, binary.LittleEndian, m.offset)
	_ = binary.Write(buf, binary.LittleEndian, m.blockSize)
	return buf.Bytes()
}
func decodeBlockMeta(r io.Reader) BlockMeta {
	var m BlockMeta
	m.lastKv = decodeKvPair(r)
	_ = binary.Read(r, binary.LittleEndian, &m.offset)
	_ = binary.Read(r, binary.LittleEndian, &m.blockSize)
	return m
}
