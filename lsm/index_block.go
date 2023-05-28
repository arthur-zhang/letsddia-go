package lsm

import (
	"bytes"
	"encoding/binary"
	"io"
)

type IndexBlock struct {
	blockMetaList []BlockMeta
	totalSize     uint64
}

func NewIndexBlock() IndexBlock {
	return IndexBlock{
		blockMetaList: make([]BlockMeta, 0),
		totalSize:     0,
	}
}
func (b *IndexBlock) appendBlockIndex(meta BlockMeta) {
	b.totalSize += meta.size()
	b.blockMetaList = append(b.blockMetaList, meta)
}
func (b *IndexBlock) encode() []byte {
	buf := new(bytes.Buffer)
	for _, meta := range b.blockMetaList {
		_ = binary.Write(buf, binary.LittleEndian, meta.encode())
	}
	return buf.Bytes()
}
func decodeIndexBlock(r io.Reader, indexSize uint64) IndexBlock {
	blockMetaList := make([]BlockMeta, 0)
	idx := uint64(0)
	for idx < indexSize {
		meta := decodeBlockMeta(r)
		blockMetaList = append(blockMetaList, meta)
		idx += meta.size()
	}

	return IndexBlock{
		blockMetaList: blockMetaList,
		totalSize:     indexSize,
	}
}
