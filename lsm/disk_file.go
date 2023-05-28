package lsm

import (
	"commons-io/file_utils"
	"errors"
	"lo"
	"os"
	"path/filepath"
)

type DiskFile struct {
	file             *os.File
	fileId           uint32
	path             string
	blockMetaList    []BlockMeta
	blockIndexOffset uint64
	blockIndexSize   uint64
	fileSize         uint64
}

func createNewDiskFile(baseDir string, fileId uint32) *DiskFile {
	datFileName := formatDatFileName(fileId)
	path := filepath.Join(baseDir, datFileName)
	file, err := file_utils.NewOpenOptions().
		Read(true).
		Write(true).
		Create(true).
		Truncate(true).
		Open(path)

	if err != nil {
		println(err.Error())
		panic(err)
	}
	return &DiskFile{
		fileId:           fileId,
		file:             file,
		path:             path,
		blockMetaList:    make([]BlockMeta, 0),
		fileSize:         0,
		blockIndexOffset: 0,
		blockIndexSize:   0,
	}
}

func openDiskFile(path string) *DiskFile {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	fid := getFileIdFromPath(path)
	diskFile := DiskFile{
		fileId:           fid,
		file:             file,
		path:             path,
		blockMetaList:    make([]BlockMeta, 0),
		fileSize:         0,
		blockIndexOffset: 0,
		blockIndexSize:   0,
	}
	err = diskFile.open()
	if err != nil {
		panic(err)
	}
	return &diskFile
}
func (f *DiskFile) open() error {
	fileLen, _ := file_utils.FileSize(f.path)
	if fileLen < META_BLOCK_SIZE {
		return errors.New("corrupted! file len too small")
	}
	metaBlockOffset := fileLen - META_BLOCK_SIZE
	_, _ = f.file.Seek(metaBlockOffset, 0)
	metaBlock := decodeMetaBlock(f.file)
	if int64(metaBlock.fileSize) != fileLen {
		return errors.New("corrupted! file size mismatch")
	}
	indexOffset := metaBlock.indexOffset
	indexSize := metaBlock.indexSize
	_, _ = f.file.Seek(int64(indexOffset), 0)
	indexBlock := decodeIndexBlock(f.file, indexSize)

	f.blockMetaList = indexBlock.blockMetaList
	f.blockIndexOffset = indexOffset
	f.blockIndexSize = indexSize
	f.fileSize = uint64(fileLen)
	return nil
}

func (f *DiskFile) NewWriter() *DiskFileWriter {
	writer := NewDiskFileWriter(f.path)
	return &writer
}

func (f *DiskFile) Iter() DiskFileIterator {
	return NewDiskFileIterator(f)
}

func (f *DiskFile) loadDataBlock(meta *BlockMeta) DataBlock {
	_, _ = f.file.Seek(int64(meta.offset), 0)
	return DecodeDataBlock(f.file, meta.blockSize)
}

func (f *DiskFile) close() {
	_ = f.file.Close()
}

type DiskFileIterator struct {
	inner            *DiskFile
	dataIndexInBlock uint64
	curBlockIndexIdx uint64
	curDataBlock     *DataBlock
	key              *Key
	value            *Value
}

func (ite *DiskFileIterator) Next() {
	ite.dataIndexInBlock += 1
	// finish current block, move to next block
	if ite.dataIndexInBlock >= uint64(len(ite.curDataBlock.KvList)) {
		ite.curBlockIndexIdx++
		ite.dataIndexInBlock = 0
		if ite.curBlockIndexIdx >= uint64(len(ite.inner.blockMetaList)) {
			ite.curDataBlock = nil
			return
		}
		meta := ite.inner.blockMetaList[ite.curBlockIndexIdx]
		dataBlock := ite.inner.loadDataBlock(&meta)
		ite.curDataBlock = &dataBlock
	}
}

func (ite *DiskFileIterator) Key() *Key {
	if ite.dataIndexInBlock >= uint64(len(ite.curDataBlock.KvList)) {
		return nil
	}
	result := ite.curDataBlock.KvList[ite.dataIndexInBlock]
	return &result.Key
}

func (ite *DiskFileIterator) Value() *Value {
	if ite.dataIndexInBlock >= uint64(len(ite.curDataBlock.KvList)) {
		return nil
	}
	result := ite.curDataBlock.KvList[ite.dataIndexInBlock]
	return &result.Value
}

func (ite *DiskFileIterator) SeekToFirst() {
	ite.curBlockIndexIdx = 0
	ite.dataIndexInBlock = 0
	blockMeta := ite.inner.blockMetaList[0]
	dataBlock := ite.inner.loadDataBlock(&blockMeta)
	ite.curDataBlock = &dataBlock
}

func (ite *DiskFileIterator) SeekToLast() {

}

func (ite *DiskFileIterator) Seek(key *Key) {
	ite.seekTo(&KVPair{Key: *key})
}

func (ite *DiskFileIterator) Valid() bool {
	if ite.curBlockIndexIdx >= uint64(len(ite.inner.blockMetaList)) {
		return false
	}
	if ite.curDataBlock == nil {
		return false
	}
	if (ite.curBlockIndexIdx == uint64(len(ite.inner.blockMetaList)-1)) && ite.dataIndexInBlock >= uint64(len(ite.curDataBlock.KvList)) {
		return false
	}
	if ite.curDataBlock == nil {
		meta := ite.inner.blockMetaList[ite.curBlockIndexIdx]
		dataBlock := ite.inner.loadDataBlock(&meta)
		ite.curDataBlock = &dataBlock
	}

	result := ite.curDataBlock.KvList[ite.dataIndexInBlock]
	ite.key = &result.Key
	ite.value = &result.Value
	return true
}

func (ite *DiskFileIterator) Close() {
	//TODO implement me
	panic("implement me")
}

func NewDiskFileIterator(inner *DiskFile) DiskFileIterator {
	file := openDiskFile(inner.path)
	return DiskFileIterator{
		inner:            file,
		dataIndexInBlock: 0,
		curBlockIndexIdx: 0,
		curDataBlock:     nil,
	}
}

func (ite *DiskFileIterator) seekTo(target *KVPair) {
	blockMetaList := ite.inner.blockMetaList
	// Find the smallest block meta which has the lastKV >= target.
	smallestBlockMeta, blockIndexIdx, found := lo.FindIndexOf(blockMetaList, func(meta BlockMeta) bool {
		return meta.lastKv.Key >= target.Key
	})
	if !found {
		ite.curBlockIndexIdx = uint64(len(blockMetaList) - 1)
		ite.dataIndexInBlock = uint64(len(ite.curDataBlock.KvList))
		ite.curDataBlock = nil
		return
	}
	ite.curBlockIndexIdx = uint64(blockIndexIdx)
	smallestBlock := ite.inner.loadDataBlock(&smallestBlockMeta)
	ite.dataIndexInBlock = 0
	_, idx, found := lo.FindIndexOf(smallestBlock.KvList, func(kv KVPair) bool {
		return kv.Key >= target.Key
	})
	if !found {
		ite.curBlockIndexIdx = uint64(len(blockMetaList) - 1)
		ite.dataIndexInBlock = uint64(len(ite.curDataBlock.KvList))
		ite.curDataBlock = nil
		return
	}

	ite.dataIndexInBlock = uint64(idx)
	ite.curDataBlock = &smallestBlock
}
