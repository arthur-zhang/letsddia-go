package lsm

import (
	"commons-io/file_utils"
	"os"
)

type DiskFileWriter struct {
	path             string
	curBlockWriter   DataBlock
	indexBlockWriter IndexBlock
	file             *os.File
	blockOffset      uint64
	indexOffset      uint64
	indexSize        uint64
	fileSize         uint64
	blockCount       uint32
}

func NewDiskFileWriter(path string) DiskFileWriter {
	file, _ := file_utils.
		NewOpenOptions().
		Write(true).
		Read(true).
		Create(true).
		Truncate(true).
		Open(path)

	return DiskFileWriter{
		path:             path,
		curBlockWriter:   NewDataBlock(),
		indexBlockWriter: NewIndexBlock(),
		file:             file,
		blockOffset:      0,
		indexOffset:      0,
		indexSize:        0,
		fileSize:         0,
		blockCount:       1,
	}
}
func (w *DiskFileWriter) rotateBlockWrite() {

	lastKv := w.curBlockWriter.LastKv()
	if lastKv == nil {
		return
	}
	blockWriter := w.curBlockWriter
	blockData := blockWriter.Encode()
	_, _ = w.file.Write(blockData)

	blockDataLen := uint64(len(blockData))
	w.fileSize += blockDataLen

	w.indexBlockWriter.appendBlockIndex(NewBlockMeta(*lastKv, w.blockOffset, blockDataLen))
	w.blockOffset += blockDataLen
	w.curBlockWriter = NewDataBlock()
	w.blockCount++

}
func (w *DiskFileWriter) Add(pair KVPair) {
	if !w.curBlockWriter.hasSpace(&pair) {
		w.rotateBlockWrite()
	}
	w.curBlockWriter.Add(&pair)
}
func (w *DiskFileWriter) AppendIndexBlock() {
	if w.curBlockWriter.kvCount() > 0 {
		w.rotateBlockWrite()
	}
	w.indexOffset = w.blockOffset
	indexBlockWriter := w.indexBlockWriter
	w.indexBlockWriter = NewIndexBlock()
	indexBlockData := indexBlockWriter.encode()
	w.blockOffset += uint64(len(indexBlockData))
	_, _ = w.file.Write(indexBlockData)
	w.fileSize += uint64(len(indexBlockData))
	w.indexSize = uint64(len(indexBlockData))
}
func (w *DiskFileWriter) AppendMetaBlock() {
	w.fileSize += META_BLOCK_SIZE
	metaBlock := NewMetaBlock(w.fileSize, w.blockCount, w.indexOffset, w.indexSize)
	metaBlockData := metaBlock.encode()
	_, _ = w.file.Write(metaBlockData)
}
func (w *DiskFileWriter) Close() {
	_ = w.file.Close()
}
