package hashmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Key int

func (key Key) hash() int {
	return int(key)
}
func (key Key) Eq(other Key) bool {
	return key == other
}

func TestHashMap(t *testing.T) {

	m := NewHashMap[Key, int]()
	for i := 0; i < 1000; i++ {
		m.Put(Key(i), i)
	}
	for i := 0; i < 1000; i++ {
		if i%3 == 0 {
			m.Delete(Key(i))
		}
	}
	for i := 0; i < 1000; i++ {
		v := m.Get(Key(i))
		if i%3 == 0 {
			assert.Nil(t, v)
		} else {
			assert.Equal(t, i, *v)
		}
	}
}
