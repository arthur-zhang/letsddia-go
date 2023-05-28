package lsm

import (
	"bytes"
	"encoding/binary"
	"io"
)

const META_BLOCK_SIZE = 8 + 4 + 8 + 8 + 8
const DISK_FILE_MAGIC = uint64(0xC8A3017F464D534C)

type MetaBlock struct {
	fileSize    uint64
	blockCount  uint32
	indexOffset uint64
	indexSize   uint64
	magic       uint64
}

func NewMetaBlock(fileSize uint64, blockCount uint32, indexOffset uint64, indexSize uint64) MetaBlock {
	return MetaBlock{
		fileSize:    fileSize,
		blockCount:  blockCount,
		indexOffset: indexOffset,
		indexSize:   indexSize,
		magic:       DISK_FILE_MAGIC,
	}
}
func (b *MetaBlock) encode() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, b.fileSize)
	_ = binary.Write(buf, binary.LittleEndian, b.blockCount)
	_ = binary.Write(buf, binary.LittleEndian, b.indexOffset)
	_ = binary.Write(buf, binary.LittleEndian, b.indexSize)
	_ = binary.Write(buf, binary.LittleEndian, b.magic)
	return buf.Bytes()
}

func decodeMetaBlock(r io.Reader) MetaBlock {
	var fileSize uint64
	var blockCount uint32
	var indexOffset uint64
	var indexSize uint64
	var magic uint64
	_ = binary.Read(r, binary.LittleEndian, &fileSize)
	_ = binary.Read(r, binary.LittleEndian, &blockCount)
	_ = binary.Read(r, binary.LittleEndian, &indexOffset)
	_ = binary.Read(r, binary.LittleEndian, &indexSize)
	_ = binary.Read(r, binary.LittleEndian, &magic)
	return MetaBlock{
		fileSize:    fileSize,
		blockCount:  blockCount,
		indexOffset: indexOffset,
		indexSize:   indexSize,
		magic:       magic,
	}
}
