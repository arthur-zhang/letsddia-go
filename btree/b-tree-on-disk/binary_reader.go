package b_tree_on_disk

import (
	"io"
	"math"
)

// BinaryReader is a binary big endian file format reader.
type BinaryReader struct {
	buf []byte
	pos uint64
	eof bool
}

// NewBinaryReader returns a big endian binary file format reader.
func NewBinaryReader(buf []byte) *BinaryReader {
	if math.MaxUint32 < uint(len(buf)) {
		return &BinaryReader{nil, 0, true}
	}
	return &BinaryReader{buf, 0, false}
}

// Seek set the reader position in the buffer.
func (r *BinaryReader) Seek(pos uint64) error {
	if uint64(len(r.buf)) < pos {
		r.eof = true
		return io.EOF
	}
	r.pos = pos
	r.eof = false
	return nil
}

func (r *BinaryReader) Skip(n uint64) error {
	if uint64(len(r.buf))-r.pos < n {
		r.eof = true
		return io.EOF
	} else {
		r.pos += n
	}
	return nil
}

// Pos returns the reader's position.
func (r *BinaryReader) Pos() uint64 {
	return r.pos
}

// Len returns the remaining length of the buffer.
func (r *BinaryReader) Len() uint64 {
	return uint64(len(r.buf)) - r.pos
}

// EOF returns true if we reached the end-of-file.
func (r *BinaryReader) EOF() bool {
	return r.eof
}

// ReadBytes reads n bytes.
func (r *BinaryReader) ReadBytes(n uint64) ([]byte, error) {
	if r.eof || uint64(len(r.buf))-r.pos < n {
		r.eof = true
		return nil, io.EOF
	}
	buf := r.buf[r.pos : r.pos+n]
	r.pos += n
	return buf, nil
}
func (r *BinaryReader) ReadBytesAt(n, offset uint64) ([]byte, error) {
	err := r.Seek(offset)
	if err != nil {
		return nil, err
	}

	return r.ReadBytes(n)
}

// ReadByte reads a single byte.
func (r *BinaryReader) ReadByte() (byte, error) {
	b, err := r.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// ReadUint8 reads a uint8.
func (r *BinaryReader) ReadUint8() (uint8, error) {
	return r.ReadByte()
}

type Endian interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}

func (r *BinaryReader) ReadUint16(endian Endian) (uint16, error) {
	b, err := r.ReadBytes(2)
	if err != nil {
		return 0, io.EOF
	}
	return endian.Uint16(b), nil
}
func (r *BinaryReader) ReadUint64(endian Endian) (uint64, error) {
	b, err := r.ReadBytes(8)
	if err != nil {
		return 0, io.EOF
	}
	return endian.Uint64(b), nil
}
func (r *BinaryReader) ReadUint64At(offset uint64, endian Endian) (uint64, error) {
	err := r.Seek(offset)
	if err != nil {
		return 0, err
	}
	return r.ReadUint64(endian)
}

func (r *BinaryReader) ReadInt64At(offset uint64, endian Endian) (int64, error) {
	readUint64, err := r.ReadUint64At(offset, endian)
	if err != nil {
		return 0, err
	}
	return int64(readUint64), nil
}

type BinaryWriter struct {
	buf [PageSize]byte
	pos uint64
}

func NewBinaryWriter(buf [PageSize]byte) BinaryWriter {
	return BinaryWriter{buf, 0}
}

func (w *BinaryWriter) WriteBool(v bool) {
	var data uint8
	if v {
		data = 1
	} else {
		data = 0
	}
	w.WriteU8(data)
}
func (w *BinaryWriter) WriteU8(v uint8) {
	w.buf[w.pos] = v
	w.pos += 1
}

func (w *BinaryWriter) WriteUint16(endian Endian) {
	endian.PutUint16(w.buf[w.pos:], uint16(w.pos))
	w.pos += 2
}

func (w *BinaryWriter) WriteBytes(bytes []byte) {
	copy(w.buf[w.pos:], bytes)
	w.pos += uint64(len(bytes))
}
func (w *BinaryWriter) WriteUint64(endian Endian, v uint64) {
	endian.PutUint64(w.buf[w.pos:], v)
	w.pos += 8
}

func (w *BinaryWriter) WriteInt64(endian Endian, v int64) {
	endian.PutUint64(w.buf[w.pos:], uint64(v))
	w.pos += 8
}
