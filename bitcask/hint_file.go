package bitcask

import (
	"commons-io/byteorder"
	"commons-io/file_utils"
	"encoding/binary"
	"os"
)

type HintFile struct {
	path string
	file *os.File
}

func OpenHintFile(path string, options *file_utils.OpenOptions) HintFile {
	file, _ := options.Open(path)
	return HintFile{
		path: path,
		file: file,
	}
}
func (f *HintFile) put(key Key, entry KeyDirEntry) {
	w := byteorder.NewBinaryIO(f.file)
	_ = w.WriteUint32(binary.BigEndian, uint32(len(key)))
	_ = w.WriteBytes([]byte(key))
	_ = w.WriteUint32(binary.BigEndian, entry.valueSz)
	_ = w.WriteUint32(binary.BigEndian, entry.valuePos)
	_ = w.WriteUint32(binary.BigEndian, entry.tstamp)
}
func (f *HintFile) Close() {
	_ = f.file.Close()
}

type HintRecord struct {
	key      Key
	valueSz  uint32
	valuePos uint32
	tstamp   uint32
}

type HintIterator struct {
	inner HintFile
	r     byteorder.BinaryIO
}

func (f *HintFile) NewIterator() *HintIterator {
	inner := OpenHintFile(f.path, file_utils.NewOpenOptions().Read(true))
	return &HintIterator{
		inner: inner,
		r:     byteorder.NewBinaryIO(inner.file)}
}
func (i *HintIterator) Next() *HintRecord {
	ksz, err := i.r.ReadUint32(binary.BigEndian)
	if err != nil {
		return nil
	}
	key, err := i.r.ReadBytes(ksz)
	if err != nil {
		return nil
	}
	valueSz, err := i.r.ReadUint32(binary.BigEndian)
	if err != nil {
		return nil
	}
	valuePos, err := i.r.ReadUint32(binary.BigEndian)
	if err != nil {
		return nil
	}
	tstamp, err := i.r.ReadUint32(binary.BigEndian)
	if err != nil {
		return nil
	}
	return &HintRecord{
		key:      Key(key),
		valueSz:  valueSz,
		valuePos: valuePos,
		tstamp:   tstamp,
	}
}
