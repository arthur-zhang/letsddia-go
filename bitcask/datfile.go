package bitcask

import (
	"commons-io/byteorder"
	"commons-io/file_utils"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
)

type DatFile struct {
	id     uint32
	path   string
	file   *os.File
	offset uint32
}

type DatFileBuilder struct {
	baseDir_     string
	fileId_      uint32
	openOptions_ *file_utils.OpenOptions
}

func newDatFileBuilder() *DatFileBuilder {
	return &DatFileBuilder{}
}
func (builder *DatFileBuilder) path(path string) *DatFileBuilder {
	builder.fileId_ = getFileIdFromPath(path)
	builder.baseDir_ = filepath.Dir(path)
	return builder
}
func (builder *DatFileBuilder) baseDir(baseDir string) *DatFileBuilder {
	builder.baseDir_ = baseDir
	return builder
}
func (builder *DatFileBuilder) fileId(fileId uint32) *DatFileBuilder {
	builder.fileId_ = fileId
	return builder
}
func (builder *DatFileBuilder) openOptions(openOptions *file_utils.OpenOptions) *DatFileBuilder {
	builder.openOptions_ = openOptions
	return builder
}
func (builder *DatFileBuilder) build() DatFile {

	filePath := filepath.Join(builder.baseDir_, formatDatFileName(builder.fileId_))
	file, err := builder.openOptions_.Open(filePath)
	if err != nil {
		panic(err)
	}
	return DatFile{
		id:     builder.fileId_,
		path:   filePath,
		file:   file,
		offset: 0,
	}
}

func (f *DatFile) Close() {
	_ = f.file.Close()
}
func (f *DatFile) getOffset() uint32 {
	return f.offset
}

func (f *DatFile) writeBlock(block *Block) {
	_, _ = f.file.Write(block.serialize())
}
func (f *DatFile) write(tstamp uint32, key Key, value Value) uint32 {
	block := NewBlock(tstamp, key, value)
	fileOffset := f.offset
	f.writeBlock(&block)
	f.offset += block.size()
	return fileOffset
}
func readBlockAt(file *os.File, pos uint32) *Block {
	r := byteorder.NewBinaryIO(file)
	err := r.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return nil
	}
	crc, err := r.ReadUint32(binary.LittleEndian)
	if err != nil {
		return nil
	}
	tstamp, err := r.ReadUint32(binary.LittleEndian)
	if err != nil {
		return nil
	}
	ksz, err := r.ReadUint32(binary.LittleEndian)
	if err != nil {
		return nil
	}
	valueSz, err := r.ReadUint32(binary.LittleEndian)
	if err != nil {
		return nil
	}
	key, err := r.ReadBytes(ksz)
	if err != nil {
		return nil
	}
	value, err := r.ReadBytes(valueSz)
	if err != nil {
		return nil
	}
	return &Block{
		crc:     crc,
		tstamp:  tstamp,
		ksz:     ksz,
		valueSz: valueSz,
		key:     Key(key),
		value:   value,
	}
}

func (f *DatFile) readValueAt(valueSz, pos uint32) (Value, error) {
	r := byteorder.NewBinaryIO(f.file)
	err := r.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return nil, err
	}
	bytes, err := r.ReadBytes(valueSz)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type BlockIterator struct {
	pos  uint32
	file *os.File
}
type IterItem struct {
	pos   uint32
	block *Block
}

func (i *BlockIterator) Next() *IterItem {
	block := readBlockAt(i.file, i.pos)
	if block == nil {
		return nil
	}
	pos := i.pos
	i.pos += block.size()
	return &IterItem{
		pos:   pos,
		block: block,
	}
}

func (f *DatFile) NewIterator() BlockIterator {
	file, err := file_utils.NewOpenOptions().Read(true).Open(f.path)
	if err != nil {
		panic(err)
	}
	return BlockIterator{
		pos:  0,
		file: file,
	}
}
