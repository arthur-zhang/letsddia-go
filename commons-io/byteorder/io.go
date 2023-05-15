package byteorder

import (
	"io"
)

type Endian interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}
type ReadBytesExt[E Endian] interface {
	ReadUint16(endian E) (uint16, error)
	ReadBytes(n int) ([]byte, error)
}
type BinaryIO struct {
	inner io.ReadWriteSeeker
}

func NewBinaryIO(rws io.ReadWriteSeeker) BinaryIO {
	return BinaryIO{inner: rws}
}
func (rws *BinaryIO) Seek(offset int64, whence int) error {
	_, err := rws.inner.Seek(offset, whence)
	return err
}
func (rws *BinaryIO) ReadBytes(n uint32) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rws.inner.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func (rws *BinaryIO) ReadUint8() (uint8, error) {
	bytes, err := rws.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}
func (rws *BinaryIO) ReadInt8() (int8, error) {
	n, err := rws.ReadUint8()
	return int8(n), err
}

func (rws *BinaryIO) ReadUint16(endian Endian) (uint16, error) {
	bytes, err := rws.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return endian.Uint16(bytes), nil
}

func (rws *BinaryIO) ReadUint32(endian Endian) (uint32, error) {
	bytes, err := rws.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return endian.Uint32(bytes), nil
}
func (rws *BinaryIO) ReadUint64(endian Endian) (uint64, error) {
	bytes, err := rws.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return endian.Uint64(bytes), nil
}
func (rws *BinaryIO) WriteBytes(data []byte) error {
	_, err := rws.inner.Write(data)
	return err
}
func (rws *BinaryIO) WriteUint8(n uint8) error {
	bytes := make([]byte, 1)
	bytes[0] = n
	return rws.WriteBytes(bytes)
}
func (rws *BinaryIO) WriteInt8(n int8) error {
	bytes := make([]byte, 1)
	bytes[0] = uint8(n)
	return rws.WriteBytes(bytes)
}
func (rws *BinaryIO) WriteUint16(endian Endian, n uint16) error {
	bytes := make([]byte, 2)
	endian.PutUint16(bytes, n)
	return rws.WriteBytes(bytes)
}
func (rws *BinaryIO) WriteInt16(endian Endian, n int16) error {
	return rws.WriteUint16(endian, uint16(n))
}
func (rws *BinaryIO) WriteUint32(endian Endian, n uint32) error {
	bytes := make([]byte, 4)
	endian.PutUint32(bytes, n)
	return rws.WriteBytes(bytes)
}
func (rws *BinaryIO) WriteInt32(endian Endian, n int32) error {
	return rws.WriteUint32(endian, uint32(n))
}

func (rws *BinaryIO) WriteUint64(endian Endian, n uint64) error {
	bytes := make([]byte, 8)
	endian.PutUint64(bytes, n)
	return rws.WriteBytes(bytes)
}

func (rws *BinaryIO) WriteInt64(endian Endian, n int64) error {
	return rws.WriteUint64(endian, uint64(n))
}
