package bitcask

import (
	"commons-io/file_utils"
	"path/filepath"
)

const RemoveTombstone = "__DEL__"

type Bitcask interface {
}
type Opts struct {
	dataFileLimit uint32
}
type Key string
type Value []byte

type KeyDirEntry struct {
	fileId   uint32
	valueSz  uint32
	valuePos uint32
	tstamp   uint32
}
type KeyDir map[Key]KeyDirEntry
type BitcaskHandle struct {
	opts          Opts
	baseDir       string
	nextFileId    uint32
	keyDir        KeyDir
	activeDatFile *DatFile
}

func (h *BitcaskHandle) createNewDatFile(fileId uint32) {
	file := newDatFileBuilder().baseDir(h.baseDir).fileId(fileId).openOptions(file_utils.NewOpenOptions().Read(true).Write(true).Create(true)).build()
	h.activeDatFile = &file
}
func (h *BitcaskHandle) getDatFilePath(fileId uint32) string {
	return filepath.Join(h.baseDir, formatDatFileName(fileId))
}
func (h *BitcaskHandle) getHintFilePath(fileId uint32) string {
	return filepath.Join(h.baseDir, formatHintFileName(fileId))
}

func (h *BitcaskHandle) checkWrite(dataLen uint32) {
	if h.activeDatFile == nil {
		h.createNewDatFile(h.nextFileId)
		return
	}

	if h.activeDatFile.getOffset()+HEADER_SIZE+dataLen > h.opts.dataFileLimit {
		h.nextFileId++
		h.createNewDatFile(h.nextFileId)
	}
}
func (h *BitcaskHandle) loadFilesInDir(datfiles []string) {
	if len(datfiles) == 0 {
		return
	}
	for _, path := range datfiles {
		file := newDatFileBuilder().path(path).openOptions(file_utils.NewOpenOptions().Read(true)).build()
		fileId := file.id
		hintFilePath := h.getHintFilePath(fileId)

		if file_utils.FileExists(hintFilePath) {
			hintFile := OpenHintFile(hintFilePath, file_utils.NewOpenOptions().Read(true))
			iter := hintFile.NewIterator()
			for item := iter.Next(); item != nil; item = iter.Next() {
				h.keyDir[item.key] = KeyDirEntry{
					fileId:   fileId,
					valueSz:  item.valueSz,
					valuePos: item.valuePos,
					tstamp:   item.tstamp,
				}
			}
		} else {
			iter := file.NewIterator()
			for item := iter.Next(); item != nil; item = iter.Next() {
				block := item.block
				if isTombstone(block.value) {
					delete(h.keyDir, block.key)
					continue
				}
				pos := item.pos
				entryInMem, found := h.keyDir[block.key]

				if !found || (entryInMem.tstamp < block.tstamp) {
					h.keyDir[block.key] = KeyDirEntry{
						fileId:   fileId,
						valueSz:  block.valueSz,
						valuePos: pos + HEADER_SIZE + block.ksz,
						tstamp:   block.tstamp,
					}
				}
			}
		}
	}
}
func Open(baseDir string, opts Opts) BitcaskHandle {
	err := file_utils.Mkdirs(baseDir)
	if err != nil {
		panic(err)
	}
	files := getDatFiles(baseDir)
	nextFileId := getNextFileId(files)

	db := BitcaskHandle{
		opts:          opts,
		baseDir:       baseDir,
		nextFileId:    nextFileId,
		keyDir:        make(map[Key]KeyDirEntry),
		activeDatFile: nil,
	}
	db.loadFilesInDir(files)
	return db
}
func isTombstone(value Value) bool {
	return string(value) == RemoveTombstone
}
func (h *BitcaskHandle) Get(key Key) (Value, error) {
	value, found := h.keyDir[key]
	if !found {
		return nil, nil
	}

	file := newDatFileBuilder().baseDir(h.baseDir).fileId(value.fileId).openOptions(file_utils.NewOpenOptions().Read(true)).build()
	defer file.Close()

	valueData, err := file.readValueAt(value.valueSz, value.valuePos)
	if err != nil {
		return nil, err
	}
	if isTombstone(valueData) {
		return nil, nil
	}
	return valueData, nil
}

func (h *BitcaskHandle) Put(key Key, value Value) {
	dataLen := len(key) + len(value)
	h.checkWrite(uint32(dataLen))
	tstamp := nowTs()
	offset := h.activeDatFile.write(tstamp, key, value)
	h.keyDir[key] = KeyDirEntry{
		fileId:   h.activeDatFile.id,
		valueSz:  uint32(len(value)),
		valuePos: offset + HEADER_SIZE + uint32(len(key)),
		tstamp:   tstamp,
	}
}
func (h *BitcaskHandle) Delete(key Key) {
	if _, found := h.keyDir[key]; !found {
		return
	}
	h.Put(key, []byte(RemoveTombstone))
	delete(h.keyDir, key)
}

func (h *BitcaskHandle) Merge() {
	files := getDatFiles(h.baseDir)
	if len(files) <= 1 {
		return
	}
	filesToMerge := files[:len(files)-1]
	for _, s := range filesToMerge {
		println("files to merge: ", s)
	}
	lastId := getFileIdFromPath(filesToMerge[len(filesToMerge)-1])
	keyDir := make(KeyDir)
	tmpDir := filepath.Join(h.baseDir, "/tmp")
	err := file_utils.Mkdirs(tmpDir)
	if err != nil {
		panic(err)
	}
	_ = file_utils.CleanDirectory(tmpDir)
	tmpFilePath := filepath.Join(tmpDir, formatDatFileName(lastId))
	//tmpFile := newDatFileBuilder().baseDir(tmpDir).fileId(lastId).flag(os.O_TRUNC | os.O_CREATE | os.O_RDWR).build()
	tmpFile := newDatFileBuilder().baseDir(tmpDir).fileId(lastId).openOptions(file_utils.NewOpenOptions().Truncate(true).Create(true).Read(true).Write(true)).build()

	tmpHintFilePath := filepath.Join(tmpDir, formatHintFileName(lastId))

	//tmpHintFile := OpenHintFile(tmpHintFilePath, false)
	tmpHintFile := OpenHintFile(tmpHintFilePath, file_utils.NewOpenOptions().Read(true).Write(true).Create(true))

	defer tmpFile.Close()
	defer tmpHintFile.Close()

	// process newest file first
	for i := len(filesToMerge) - 1; i >= 0; i-- {
		file := filesToMerge[i]
		fileId := getFileIdFromPath(file)
		println("merge file: ", file)
		datFile := newDatFileBuilder().baseDir(h.baseDir).fileId(fileId).openOptions(file_utils.NewOpenOptions().Read(true)).build()
		defer datFile.Close()
		iter := datFile.NewIterator()
		for item := iter.Next(); item != nil; item = iter.Next() {
			block := item.block
			if _, ok := keyDir[block.key]; ok {
				continue
			}
			offset := tmpFile.write(block.tstamp, block.key, block.value)

			keyDirEntry := KeyDirEntry{
				fileId:   lastId, // set to lastId
				valueSz:  block.valueSz,
				valuePos: offset + HEADER_SIZE + block.ksz,
				tstamp:   block.tstamp,
			}
			tmpHintFile.put(block.key, keyDirEntry)
			keyDir[block.key] = keyDirEntry
		}
	}

	err = file_utils.CopyFile(tmpFilePath, h.getDatFilePath(lastId))
	if err != nil {
		panic(err)
	}
	err = file_utils.CopyFile(tmpHintFilePath, h.getHintFilePath(lastId))

	if err != nil {
		panic(err)
	}
	filesToDelete := filesToMerge[:len(filesToMerge)-1]
	for _, file := range filesToDelete {
		_ = file_utils.Delete(file)
		fileId := getFileIdFromPath(file)
		_ = file_utils.Delete(h.getHintFilePath(fileId))
	}
	_ = file_utils.DeleteDirectory(tmpDir)
}
