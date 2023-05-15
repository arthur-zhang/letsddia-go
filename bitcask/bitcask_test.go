package bitcask

import (
	"commons-io/file_utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	bitcask := BitcaskHandle{
		opts:          Opts{dataFileLimit: 128},
		baseDir:       "/tmp/bitcask",
		nextFileId:    0,
		keyDir:        make(map[Key]KeyDirEntry),
		activeDatFile: nil,
	}
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("hello#%d", i)
		bitcask.Put(Key(key), []byte("world"))
	}

}

func TestOpen(t *testing.T) {
	db := Open("/tmp/bitcask", Opts{dataFileLimit: 128})

	for key, entry := range db.keyDir {
		fmt.Printf("%s: %v\n", key, entry)
	}
}

func TestGet(t *testing.T) {
	err := file_utils.CleanDirectory("/tmp/bitcask")
	if err != nil {
		t.Error(err)
	}
	db := Open("/tmp/bitcask", Opts{dataFileLimit: 128})

	for i := 0; i < 10; i++ {
		if (i % 2) == 0 {
			db.Put(Key(fmt.Sprintf("hello#%d", i)), Value(fmt.Sprintf("world#%d", i)))
		}
	}
	for i := 0; i < 10; i++ {
		value, err := db.Get(Key(fmt.Sprintf("hello#%d", i)))
		if err != nil {
			t.Error(err)
		}

		if (i % 2) == 0 {
			assert.Equal(t, Value(fmt.Sprintf("world#%d", i)), value)
		} else {
			assert.True(t, value == nil)
		}
	}
}

func TestDelete(t *testing.T) {
	err := file_utils.CleanDirectory("/tmp/bitcask")
	if err != nil {
		t.Error(err)
	}
	db := Open("/tmp/bitcask", Opts{dataFileLimit: 128})

	for i := 0; i < 10; i++ {
		db.Put(Key(fmt.Sprintf("hello#%d", i)), Value(fmt.Sprintf("world#%d", i)))
	}
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			db.Delete(Key(fmt.Sprintf("hello#%d", i)))
		}
	}

	for i := 0; i < 10; i++ {
		value, err := db.Get(Key(fmt.Sprintf("hello#%d", i)))
		if err != nil {
			t.Error(err)
		}

		if (i % 2) == 0 {
			assert.True(t, value == nil)
		} else {
			assert.Equal(t, Value(fmt.Sprintf("world#%d", i)), value)
		}
	}
}

func TestMerge(t *testing.T) {
	db := Open("/tmp/bitcask", Opts{dataFileLimit: 20})

	for i := 0; i < 100; i++ {
		db.Put(Key(fmt.Sprintf("hello#%d", i)), Value(fmt.Sprintf("world#%d", i)))
	}
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			db.Delete(Key(fmt.Sprintf("hello#%d", i)))
		}
	}

	db.Merge()
}

func TestMergeThenGet(t *testing.T) {
	{

		db := Open("/tmp/bitcask", Opts{dataFileLimit: 20})

		for i := 0; i < 100; i++ {
			db.Put(Key(fmt.Sprintf("hello#%d", i)), Value(fmt.Sprintf("world#%d", i)))
		}
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				db.Delete(Key(fmt.Sprintf("hello#%d", i)))
			}
		}

		db.Merge()
	}
	{
		db := Open("/tmp/bitcask", Opts{dataFileLimit: 20})
		for i := 0; i < 100; i++ {
			v, err := db.Get(Key(fmt.Sprintf("hello#%d", i)))
			if err != nil {
				t.Error(err)
			}

			if i%2 == 0 {
				assert.True(t, v == nil, fmt.Sprintf("hello#%d", i))
			} else {
				assert.Equal(t, Value(fmt.Sprintf("world#%d", i)), v)
			}
		}
	}
}
