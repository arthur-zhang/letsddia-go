package skiplist

import (
	"fmt"
	"math/rand"

	"golang.org/x/exp/constraints"
)

type Node[K constraints.Ordered, V any] struct {
	key     *K
	value   *V
	forward []*Node[K, V]
}

func (n *Node[K, V]) hasNext(level int) bool {
	return n.next(level) != nil
}
func (n *Node[K, V]) next(level int) *Node[K, V] {
	return n.forward[level]
}

func NewNode[K constraints.Ordered, V any](level int, key *K, value *V) *Node[K, V] {
	return &Node[K, V]{
		key:     key,
		value:   value,
		forward: make([]*Node[K, V], level),
	}
}

type SkipList[K constraints.Ordered, V any] struct {
	Head  *Node[K, V]
	Level int
	size  int
}

// NewSkipList The new skip list is initialized :
// 1. the level of the skiplist is 1
// 2. the header of the skiplist has forward-pointers at level one through MaxLevel
// 3. all forward pointers of head point to nil
func NewSkipList[K constraints.Ordered, V any]() *SkipList[K, V] {
	var zeroK K
	var zeroV V
	return &SkipList[K, V]{
		Head:  NewNode[K, V](MaxLevel, &zeroK, &zeroV),
		Level: 1,
	}
}

const MaxLevel = 32
const P = 4

func randomLevel() int {
	level := 1
	for rand.Float32() < 1.0/float32(P) && level < MaxLevel {
		level += 1
	}
	return level
}
func (list *SkipList[K, V]) Display() {
	for i := list.Level - 1; i >= 0; i-- {
		fmt.Printf("Level %d: ", i)
		x := list.Head

		for x.forward[i] != nil {
			fmt.Printf("%+v -> ", x.forward[i].key)
			x = x.forward[i]
		}
		fmt.Println("nil")
	}
}
func (list *SkipList[K, V]) Contains(key *K) bool {
	x := list.FindGreaterThanOrEqual(key)
	return x != nil && *x.key == *key
}
func (list *SkipList[K, V]) Insert(key *K, value *V) {
	// Create an update slice, which is used to store the nodes that need to update pointers in each layer,
	// and update the pointers of the update nodes according to the level of the inserted node later.
	update := make([]*Node[K, V], MaxLevel)

	x := list.Head
	// Traverse the skiplist from the top level until the position where the new key should be inserted is found.

	for i := list.Level - 1; i >= 0; i-- {
		// Traverse the current layer until a node is found whose successor node is nil or greater than or equal to the target value
		for x.forward[i] != nil && *x.forward[i].key < *key {
			x = x.forward[i]
		}
		// x.value < value <= x.forward[i].value
		// Store the current node in the update slice for later pointer updates.
		update[i] = x
	}
	// If there is a duplicate key, return directly without insertion
	x = x.forward[0]

	if x != nil && x.key == key {
		// found the key, update the value
		x.value = value
		return
	}

	// Generate a random level for the new node
	newLevel := randomLevel()
	// If the new node's level is greater than the current skiplist level, update the skiplist level and store the new layer's head node in update.
	if newLevel > list.Level {
		for i := list.Level; i < newLevel; i++ {
			update[i] = list.Head
		}
		list.Level = newLevel
	}
	// Create a new node
	x = NewNode[K, V](newLevel, key, value)
	// Insert the new node into the skiplist, similar to inserting a linked list
	for i := 0; i < newLevel; i++ {
		// Set the new node's successor node to update[i]'s successor node
		x.forward[i] = update[i].forward[i]
		// Set update[i]'s successor node to the new node
		update[i].forward[i] = x
	}
	list.size += 1
}
func (list *SkipList[K, V]) Delete(searchKey *K) {
	// Create an update slice, which is used to store the nodes that need to update pointers in each layer.
	update := make([]*Node[K, V], MaxLevel)
	x := list.Head
	// loop invariant: x.value < searchKey
	// Traverse the skiplist from the bottom layer until the node to be deleted is found.
	for i := list.Level - 1; i >= 0; i-- {
		// Traverse the current layer until a node is found whose successor node is nil or greater than or equal to the target value
		for x.forward[i] != nil && *x.forward[i].key < *searchKey {
			x = x.forward[i]
		}
		// x.value < searchKey <= x.forward[i].value
		// Store the current node in the update slice for later pointer updates.
		update[i] = x
	}
	// Move to the node to be deleted
	x = x.forward[0]

	// If the node to be deleted does not exist, return directly
	if x == nil || x.key != searchKey {
		return
	}

	// Traverse the levels and update the corresponding nodes
	for i := 0; i < list.Level; i++ {
		// If update[i]'s successor node is not the node to be deleted, skip this layer.
		if update[i].forward[i] != x {
			break
		}
		// Update update[i]'s successor node to the node to be deleted's successor node
		update[i].forward[i] = x.forward[i]
	}
	// If the deleted node is the highest level node, update the skiplist level.
	for list.Level > 1 && list.Head.forward[list.Level-1] == nil {
		list.Level--
	}
	list.size -= 1
}

func (list *SkipList[K, V]) Size() int {
	return list.size
}

func (list *SkipList[K, V]) IsEmpty() bool {
	return list.size == 0
}

func (list *SkipList[K, V]) Iterator() Iterator[K, V] {
	return &iter[K, V]{
		node: nil,
		list: list,
	}
}

func (list *SkipList[K, V]) FindLessThan(searchKey *K) *Node[K, V] {
	var result *Node[K, V]
	x := list.Head
	for i := list.Level - 1; i >= 0; i-- {
		for x.forward[i] != nil && *x.forward[i].key < *searchKey {
			x = x.forward[i]
			if i == 0 {
				result = x
			}
		}

	}
	return result
}

// FindGreaterThanOrEqual returns the first node greater than or equal to searchKey
func (list *SkipList[K, V]) FindGreaterThanOrEqual(searchKey *K) *Node[K, V] {
	x := list.Head
	for i := list.Level - 1; i >= 0; i-- {
		for x.forward[i] != nil && *x.forward[i].key < *searchKey {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	return x
}

// FindLast returns the last node in the skiplist
func (list *SkipList[K, V]) FindLast() *Node[K, V] {
	x := list.Head
	// Traverse the skiplist from the top layer
	for i := list.Level - 1; i >= 0; i-- {
		// Traverse the current layer until a node is found whose successor node is nil
		for x.hasNext(i) {
			x = x.next(i)
		}
	}
	// If the found node is the header node, return nil, indicating that the skiplist is empty.
	if x == list.Head {
		return nil
	}
	return x
}

func (list *SkipList[K, V]) Search(searchKey *K) *Node[K, V] {
	x := list.Head
	// loop invariant: x.value < searchKey
	for i := list.Level - 1; i >= 0; i-- {
		for *x.forward[i].key < *searchKey {
			x = x.forward[i]
		}
	}
	// x.value < searchKey <= x.forward[i].value
	x = x.forward[0]
	if x.key == searchKey {
		return x
	}
	return nil
}

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

type iter[K constraints.Ordered, V any] struct {
	node *Node[K, V]
	list *SkipList[K, V]
}

func (i *iter[K, V]) Next() {
	if i.node != nil {
		i.node = i.node.next(0)
	}
}

func (i *iter[K, V]) Key() (key *K) {
	return i.node.key
}

func (i *iter[K, V]) Value() (value *V) {
	return i.node.value
}

func (i *iter[K, V]) Seek(key *K) {
	i.node = i.list.FindGreaterThanOrEqual(key)
}
func (i *iter[K, V]) SeekToFirst() {
	i.node = i.list.Head.next(0)
}
func (i *iter[K, V]) SeekToLast() {
	i.node = i.list.FindLast()
}
func (i *iter[K, V]) Valid() bool {
	return i.node != nil
}
func (i *iter[K, V]) Close() {
	i.node = nil
	i.list = nil
}
