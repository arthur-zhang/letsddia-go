package lsm

import (
	"commons-io/file_utils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func formatDatFileName(fileId uint32) string {
	return fmt.Sprintf("%010d.dat", fileId)
}
func getFileIdFromPath(path string) uint32 {
	fileName := getFileNameWithoutExtension(path)
	return parseUint32(fileName)
}

func getFileNameWithoutExtension(path string) string {
	filename := filepath.Base(path)
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

func parseUint32(s string) uint32 {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(num)
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
