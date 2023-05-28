package lsm

import (
	"bytes"
	"skiplist"
)

type Key string
type Value []byte
type MemStore struct {
	kvMap     *skiplist.SkipList[Key, Value]
	dataSize  uint64
	diskStore *DiskStore
	config    *Config
}

func NewMemStore(diskStore *DiskStore, config *Config) MemStore {
	return MemStore{
		kvMap:     skiplist.NewSkipList[Key, Value](),
		dataSize:  0,
		diskStore: diskStore,
		config:    config,
	}
}
func (m *MemStore) flushIfNeed() {
	if m.dataSize >= m.config.maxMemStoreSize {
		m.flush()
	}
}
func (m *MemStore) Iter() Iterator[Key, Value] {
	return m.kvMap.Iterator()
}

func (m *MemStore) Add(seqId uint64, op uint8, key Key, value Value) {
	kv := KVPair{
		SeqId: seqId,
		Op:    op,
		Key:   key,
		Value: value,
	}
	kvBytes := kv.encode()
	m.kvMap.Insert(&key, (*Value)(&kvBytes))
	// todo add update logic
	m.dataSize += kv.size()
	m.flushIfNeed()
}

func (m *MemStore) flush() {
	if m.kvMap.IsEmpty() {
		return
	}
	kvMap := m.kvMap
	diskFile := m.diskStore.genDiskFile()
	writer := diskFile.NewWriter()

	iter := kvMap.Iterator()
	iter.SeekToFirst()
	for iter.Valid() {
		kvBytes := iter.Value()
		kv := decodeKvPair(bytes.NewBuffer(*kvBytes))
		writer.Add(kv)
		iter.Next()
	}
	writer.AppendIndexBlock()
	writer.AppendMetaBlock()
	writer.Close()
	m.dataSize = 0
	m.kvMap = skiplist.NewSkipList[Key, Value]()
}
