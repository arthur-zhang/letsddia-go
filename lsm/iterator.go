package lsm

import (
	"container/heap"
	"golang.org/x/exp/constraints"
)

type Iterator[K constraints.Ordered, V any] interface {
	Next()
	Key() *K
	Value() *V
	SeekToFirst()
	SeekToLast()
	Seek(key *K)
	Valid() bool
	Close()
}

type MultiIterator struct {
	pq       MinHeap
	inner    []Iterator[Key, Value]
	curEntry *PqEntry
}

func NewMultiIterator(inner []Iterator[Key, Value]) MultiIterator {
	pq := MinHeap{}
	heap.Init(&pq)

	return MultiIterator{
		pq:    pq,
		inner: inner,
	}
}

func (m *MultiIterator) Next() {
	if len(m.pq) == 0 {
		m.curEntry = nil
		return
	}
	entry := heap.Pop(&m.pq).(PqEntry)
	m.curEntry = &entry
	entry.iter.Next()
	if entry.iter.Valid() {
		key := entry.iter.Key()
		value := entry.iter.Value()
		heap.Push(&m.pq, PqEntry{
			key:   key,
			value: value,
			iter:  entry.iter,
		})
	}
}

func (m *MultiIterator) Key() *Key {
	return m.curEntry.key
}

func (m *MultiIterator) Value() *Value {
	return m.curEntry.value
}

func (m *MultiIterator) SeekToFirst() {
	clearPq(&m.pq)
	for _, iter := range m.inner {
		iter.SeekToFirst()
	}
	for _, iter := range m.inner {
		if iter.Valid() {
			key := iter.Key()
			value := iter.Value()
			heap.Push(&m.pq, PqEntry{key: key, value: value, iter: iter})
		}
	}
	m.Next()
}

func (m *MultiIterator) SeekToLast() {
	panic("implement me")
}

func clearPq(pq *MinHeap) {
	for len(*pq) > 0 {
		heap.Pop(pq)
	}
}
func (m *MultiIterator) Seek(key *Key) {
	clearPq(&m.pq)
	for i, iter := range m.inner {
		_ = i
		iter.Seek(key)
		if iter.Valid() {
			k := iter.Key()
			v := iter.Value()
			heap.Push(&m.pq, PqEntry{key: k, value: v, iter: iter})
		}
	}
	m.Next()
}

func (m *MultiIterator) Valid() bool {
	return len(m.pq) > 0
}

func (m *MultiIterator) Close() {
	//TODO implement me
	panic("implement me")
}

type Entry[K constraints.Ordered] struct {
	value   K
	iterIdx int
}

type PqEntry struct {
	key   *Key
	value *Value
	iter  Iterator[Key, Value]
}
type MinHeap []PqEntry

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return *h[i].key < *h[j].key }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(PqEntry))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
