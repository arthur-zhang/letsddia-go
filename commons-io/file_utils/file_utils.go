package file_utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type FileFilter func(file os.FileInfo) bool

func ByteCountToDisplaySize() {}
func Checksum()               {}
func ChecksumCRC32()          {}
func CleanDirectory(directory string) error {
	files, err := ListFiles(directory, nil)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := ForceDelete(file)
		if err != nil {
			return err
		}
	}
	return nil
}
func CleanDirectoryOnExit()             {}
func ContentEquals()                    {}
func ContentEqualsIgnoreEOL()           {}
func ConvertFileCollectionToFileArray() {}
func CopyDirectory()                    {}
func CopyDirectoryToDirectory()         {}
func CopyFile(srcFile, destFile string) error {
	err := requireFileCopy(srcFile, destFile)
	if err != nil {
		return err
	}
	err = requireFile(srcFile, "srcFile")
	if err != nil {
		return err
	}
	err = requireCanonicalPathsNotEquals(srcFile, destFile)
	if err != nil {
		return err
	}
	err = createParentDirectories(destFile)
	if err != nil {
		return err
	}
	err = requireFileIfExists(destFile, "destFile")
	if err != nil {
		return err
	}
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(destFile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	return err
}
func CopyFileToDirectory()   {}
func CopyInputStreamToFile() {}
func CopyToDirectory()       {}
func CopyToFile()            {}
func CopyURLToFile()         {}
func createParentDirectories(file string) error {
	parentFile := GetParentFile(file)
	err := Mkdirs(parentFile)
	if err != nil {
		return err
	}
	return nil
}
func Current() {}
func Delete(file string) error {
	return os.Remove(file)
}
func DeleteDirectory(directory string) error {
	if !FileExists(directory) {
		return nil
	}
	err := CleanDirectory(directory)
	if err != nil {
		return err
	}
	return Delete(directory)
}
func DeleteQuietly(file string) bool {
	isDir, err := IsDirectory(file)
	if err != nil {
		return false
	}
	if isDir {
		err = CleanDirectory(file)
		if err != nil {
			return false
		}
	}
	err = Delete(file)
	if err != nil {
		return false
	}
	return true
}
func DirectoryContains() {}
func DoCopyDirectory()   {}
func ForceDelete(file string) error {
	return os.RemoveAll(file)
}
func ForceDeleteOnExit() {}
func ForceMkdir()        {}
func ForceMkdirParent()  {}
func GetFile()           {}
func GetParentFile(file string) string {
	return filepath.Dir(file)
}
func GetTempDirectory()      {}
func GetTempDirectoryPath()  {}
func GetUserDirectory()      {}
func GetUserDirectoryPath()  {}
func IsEmptyDirectory()      {}
func IsFileNewer()           {}
func IsFileOlder()           {}
func IsRegularFile()         {}
func IsSymlink()             {}
func IterateFiles()          {}
func IterateFilesAndDirs()   {}
func LastModified()          {}
func LastModifiedFileTime()  {}
func LastModifiedUnchecked() {}
func LineIterator()          {}
func ListAccumulate()        {}

var AllTrueFilter = func(file os.FileInfo) bool {
	return true
}

var AllFileFilter = func(file os.FileInfo) bool {
	return file.Mode().IsRegular()
}
var AllDirFilter = func(file os.FileInfo) bool {
	return file.Mode().IsDir()
}

func ListFiles(directory string, fileFilter FileFilter) ([]string, error) {
	err := requireDirectoryExists(directory, "directory")
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	if fileFilter == nil {
		fileFilter = AllTrueFilter
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if fileFilter(info) {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}
	return files, nil
}
func ListFilesAndDirs() {}

func Mkdirs(directory string) error {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}
	return nil
}
func MoveDirectory()            {}
func MoveDirectoryToDirectory() {}

func MoveFile(srcFile, destFile string) error {
	err := ValidateMoveParameters(srcFile, destFile)
	if err != nil {
		return err
	}
	err = requireFile(srcFile, "srcFile")
	if err != nil {
		return err
	}
	err = requireAbsent(destFile, "destFile")
	if err != nil {
		return err
	}
	err = os.Rename(srcFile, destFile)
	if err != nil {
		return err
	}
	return nil
}

func MoveFileToDirectory(srcFile, destDir string, createDestDir bool) error {
	err := ValidateMoveParameters(srcFile, destDir)
	if err != nil {
		return err
	}
	if !FileExists(destDir) && createDestDir {
		err = Mkdirs(destDir)
		if err != nil {
			return err
		}
	}
	err = requireExistsChecked(destDir, "destDir")
	if err != nil {
		return err
	}
	err = requireDirectory(destDir, "destDir")
	if err != nil {
		return err
	}
	destFile := destDir + "/" + GetFileName(srcFile)
	err = MoveFile(srcFile, destFile)
	if err != nil {
		return err
	}
	return nil
}
func MoveToDirectory()     {}
func NewOutputStream()     {}
func OpenInputStream()     {}
func OpenOutputStream()    {}
func ReadFileToByteArray() {}
func ReadFileToString()    {}
func ReadLines()           {}
func requireAbsent(file, name string) error {
	if FileExists(file) {
		return errors.New("File system element for parameter '" + name + "' already exists: '" + file + "'")
	}
	return nil
}
func requireCanonicalPathsNotEquals(file1, file2 string) error {
	if file1 == file2 {
		return errors.New("File canonical paths must not be equal., file1=" + file1 + ", file2=" + file2)
	}
	return nil
}
func RequireCanWrite() {}
func requireDirectory(directory, name string) error {
	ok, err := IsDirectory(directory)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Parameter '" + name + "' is not a directory: '" + directory + "'")
	}
	return nil
}
func requireDirectoryExists(directory, name string) error {

	err := RequireExists(directory, name)
	if err != nil {
		return err
	}
	return requireDirectory(directory, name)
}
func RequireDirectoryIfExists() {}
func RequireExists(file, fileParamName string) error {
	if !FileExists(file) {
		return errors.New("File system element for parameter '" + fileParamName + "' does not exist: '" + file + "'")
	}
	return nil
}

func requireExistsChecked(file, fileParamName string) error {
	if !FileExists(file) {
		return errors.New("File system element for parameter '" + fileParamName + "' does not exist: '" + file + "'")
	}
	return nil
}
func requireFile(file, name string) error {
	isFile, err := IsFile(file)
	if err != nil {
		return err
	}
	if !isFile {
		return errors.New("Parameter '" + name + "' is not a file: '" + file + "'")
	}
	return nil
}
func requireFileCopy(source, destination string) error {
	return requireExistsChecked(source, "source")
}
func requireFileIfExists(file, name string) error {
	if FileExists(file) {
		return requireFile(file, name)
	}
	return nil
}
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
func SetTimes()                    {}
func SizeOf()                      {}
func SizeOfAsBigInteger()          {}
func SizeOfDirectory()             {}
func SizeOfDirectoryAsBigInteger() {}
func StreamFiles()                 {}
func ToFile()                      {}
func ToFiles()                     {}
func ToList()                      {}
func ToMaxDepth()                  {}
func ToSuffixes()                  {}
func Touch()                       {}
func ToURLs()                      {}

func ValidateMoveParameters(source, dest string) error {
	if !FileExists(source) {
		return errors.New("Source '" + source + "' does not exist")
	}
	return nil
}
func WaitFor()              {}
func Write()                {}
func WriteByteArrayToFile() {}
func WriteLines()           {}
func WriteStringToFile()    {}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// IsDirectory when returns true if the path is a dir
// when return false
//  1. if error is nil, path is not exist
//  2. if error is not nil, path is not a dir
func IsDirectory(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.Mode().IsDir(), nil
}

// IsFile when returns true if the path is a file
// when return false
//  1. if error is nil, path is not exist
//  2. if error is not nil, path is not a file
func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.Mode().IsRegular(), nil
}

func GetFileName(path string) string {
	return filepath.Base(path)
}

func Join(paths ...string) string {
	return filepath.Join(paths...)
}
