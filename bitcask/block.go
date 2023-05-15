package bitcask

import "hash/crc32"

const HEADER_SIZE = 16

type Block struct {
	crc     uint32
	tstamp  uint32
	ksz     uint32
	valueSz uint32
	key     Key
	value   Value
}

func NewBlock(tstamp uint32, key Key, value Value) Block {
	b := Block{
		tstamp:  tstamp,
		ksz:     uint32(len(key)),
		key:     key,
		valueSz: uint32(len(value)),
		value:   value,
	}
	b.crc = b.blockCrc()
	return b
}
func (b *Block) size() uint32 {
	return HEADER_SIZE + b.ksz + b.valueSz
}
func (b *Block) blockCrc() uint32 {
	crc := uint32(0)
	crc = crc32.Update(crc, crc32.IEEETable, uint32ToBytes(b.tstamp))
	crc = crc32.Update(crc, crc32.IEEETable, uint32ToBytes(b.ksz))
	crc = crc32.Update(crc, crc32.IEEETable, uint32ToBytes(b.valueSz))
	crc = crc32.Update(crc, crc32.IEEETable, []byte(b.key))
	crc = crc32.Update(crc, crc32.IEEETable, b.value)
	return crc
}
func (b *Block) serialize() []byte {
	buf := make([]byte, b.size())
	copy(buf, uint32ToBytes(b.crc))
	copy(buf[4:], uint32ToBytes(b.tstamp))
	copy(buf[8:], uint32ToBytes(b.ksz))
	copy(buf[12:], uint32ToBytes(b.valueSz))
	copy(buf[16:], b.key)
	copy(buf[16+b.ksz:], b.value)
	return buf
}
