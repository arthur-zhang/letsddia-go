package lsm

import "lo"

type DiskStore struct {
	baseDir      string
	diskFileList []*DiskFile
	nextFileId   uint32
}

func NewDiskStore(baseDir string) DiskStore {
	files := getDatFiles(baseDir)
	diskFiles := lo.Map(files, func(path string, _ int) *DiskFile {
		file := openDiskFile(path)
		return file
	})
	var newestFileId uint32 = 0
	if len(diskFiles) > 0 {
		newestFileId = diskFiles[len(diskFiles)-1].fileId
	}
	return DiskStore{
		baseDir:      baseDir,
		diskFileList: diskFiles,
		nextFileId:   newestFileId + 1,
	}
}
func (w *DiskStore) genDiskFile() *DiskFile {
	diskFile := createNewDiskFile(w.baseDir, w.nextFileId)
	w.AddDiskFile(diskFile)
	w.nextFileId += 1
	return diskFile
}

func (w *DiskStore) AddDiskFile(file *DiskFile) {
	w.diskFileList = append(w.diskFileList, file)
}

func (w *DiskStore) Close() {
	for _, file := range w.diskFileList {
		file.close()
	}
}
func (w *DiskStore) Iter() *MultiIterator {
	iters := lo.Map(w.diskFileList, func(file *DiskFile, _ int) Iterator[Key, Value] {
		iter := file.Iter()
		iter.SeekToFirst()
		return &iter
	})
	iter := NewMultiIterator(iters)
	return &iter
}
