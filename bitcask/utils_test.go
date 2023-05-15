package bitcask

import "testing"

func TestParseFileId(t *testing.T) {
	s := "0000000123"

	println(parseUint32(s))
}
func TestGetDatfiles(t *testing.T) {
	files := getDatFiles("/tmp/bitcask")
	for _, file := range files {
		println(file)
	}
}
