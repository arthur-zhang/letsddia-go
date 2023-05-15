package file_utils

import "testing"

func TestMove(t *testing.T) {
	newFile := "/tmp/bitcask/hint_file"

	err := MoveFileToDirectory("/tmp/hint_file", "/tmp/bitcask", false)
	if err != nil {
		t.Errorf("move file failed: %s", err)
		return
	}
	if !FileExists(newFile) {
		t.Errorf("file %s should exist", newFile)
	}
}

func TestCleanDirectory(t *testing.T) {
	err := CleanDirectory("/tmp/bitcask_bk")
	if err != nil {
		t.Errorf("clean directory failed: %s", err)
		return
	}
}

func TestParent(t *testing.T) {
	parent := GetParentFile("/tmp/bitcask_bk")
	println(parent)
	if parent != "/tmp" {
		t.Errorf("parent should be /tmp")
		return
	}
}
