package bitcask

import (
	"commons-io/file_utils"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func formatDatFileName(fileId uint32) string {
	return fmt.Sprintf("%010d.dat", fileId)
}
func formatHintFileName(fileId uint32) string {
	return fmt.Sprintf("%010d.hint", fileId)
}
func getFileIdFromPath(path string) uint32 {
	fileName := getFileNameWithoutExtension(path)
	return parseUint32(fileName)
}
func parseUint32(s string) uint32 {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(num)
}
func getFileNameWithoutExtension(path string) string {
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

func nowTs() uint32 {
	return uint32(time.Now().Unix())
}
func uint32ToBytes(n uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, n)
	return bytes
}

func getDatFiles(dir string) []string {
	files, err := file_utils.ListFiles(dir, func(file os.FileInfo) bool {
		return strings.HasSuffix(file.Name(), ".dat")
	})

	if err != nil {
		panic(err)
	}
	sort.Strings(files)
	return files
}
func getNextFileId(datFiles []string) uint32 {
	if len(datFiles) == 0 {
		return 0
	}
	lastFile := datFiles[len(datFiles)-1]
	return getFileIdFromPath(lastFile) + 1
}
