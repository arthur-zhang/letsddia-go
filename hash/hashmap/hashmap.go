package hashmap

import (
	"container/list"
)

type Eq[T any] interface {
	Eq(other T) bool
}
type Hasher interface {
	hash() int
}
type HashItem[T any] interface {
	Eq[T]
	Hasher
}
type KvItem[K HashItem[K], V any] struct {
	key   K
	value V
}
type HashMap[K HashItem[K], V any] struct {
	blocks []*list.List
}

const DefaultBase = 769

func NewHashMap[K HashItem[K], V any]() HashMap[K, V] {
	blocks := make([]*list.List, DefaultBase)
	for i := 0; i < DefaultBase; i++ {
		blocks[i] = list.New()
	}
	return HashMap[K, V]{
		blocks: blocks,
	}
}

func (m *HashMap[K, V]) Put(key K, value V) {
	keyHash := key.hash() % DefaultBase

	l := m.blocks[keyHash]
	for e := l.Front(); e != nil; e = e.Next() {
		item := e.Value.(KvItem[K, V])
		if item.key.Eq(key) {
			return
		}
	}
	l.PushBack(KvItem[K, V]{key: key, value: value})
}
func (m *HashMap[K, V]) Get(key K) *V {
	keyHash := key.hash() % DefaultBase

	l := m.blocks[keyHash]
	for e := l.Front(); e != nil; e = e.Next() {
		item := e.Value.(KvItem[K, V])
		if item.key.Eq(key) {
			return &item.value
		}
	}
	return nil
}

func (m *HashMap[K, V]) Delete(key K) {
	keyHash := key.hash() % DefaultBase
	l := m.blocks[keyHash]
	for e := l.Front(); e != nil; e = e.Next() {
		item := e.Value.(KvItem[K, V])
		if item.key.Eq(key) {
			l.Remove(e)
			return
		}
	}
}
