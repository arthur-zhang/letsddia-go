package b_tree_on_disk

import (
	"fmt"
	"golang.org/x/exp/slices"
)

type Btree struct {
	pager      *Pager
	rootOffset int64
	t          int
}

const N = 10

type Item struct {
	key   [N]byte
	value [N]byte
}

func (i *Item) Compare(other *Item) int {
	return slices.Compare(i.key[:], other.key[:])
}

type Comparable[T any] interface {
	Compare(other T) int
}

type Node struct {
	offset   int64
	isLeaf   bool
	items    []*Item
	children []int64
}

func NewNode(isLeaf bool) *Node {
	return &Node{
		offset:   0,
		isLeaf:   isLeaf,
		items:    make([]*Item, 0),
		children: make([]int64, 0),
	}
}
func NewBtree(path string, t int) *Btree {
	if t < 2 {
		panic(fmt.Sprintf("t must be greater than 1, got %d", t))
	}
	pager := NewPager(path)
	rootNode := NewNode(true)
	rootPage := rootNode.ToPage()
	rootPageOffset, _ := pager.WritePage(rootPage)

	return &Btree{
		pager:      pager,
		rootOffset: rootPageOffset,
		t:          t,
	}
}
func (tree *Btree) readNodeAtOffset(offset int64) *Node {
	node := tree.pager.ReadPage(offset).ToNode()
	node.offset = offset
	return node
}

func (tree *Btree) writeNodeAtOffset(node *Node) {
	tree.pager.WritePageAtOffset(node.ToPage(), node.offset)
}

func (tree *Btree) writeNewNode(node *Node) {
	offset, err := tree.pager.WritePage(node.ToPage())
	if err != nil {
		panic(err)
	}
	node.offset = offset
}

func (tree *Btree) Insert(key *Item) {
	rootOffset := tree.rootOffset
	rootNode := tree.readNodeAtOffset(rootOffset)
	// if root is full, split root
	if len(rootNode.items) == 2*tree.t-1 {
		s := NewNode(false)
		s.children = append(s.children, rootOffset)
		// must call before split to get valid file offset
		tree.writeNewNode(s)
		tree.SplitChild(s, rootNode, 0)
		tree.InsertNonFull(s, key)
		tree.rootOffset = s.offset
	} else {
		tree.InsertNonFull(rootNode, key)
		tree.writeNodeAtOffset(rootNode)
	}

}

// SplitChild split child node y of x at i
//
//	 x
//	/ \
//
// y   z
func (tree *Btree) SplitChild(x *Node, y *Node, i int) {
	T := tree.t
	z := NewNode(y.isLeaf)
	z.items = append(z.items, y.items[T:]...)
	if !z.isLeaf {
		z.children = append(z.children, y.children[T:]...)
		y.children = y.children[:T]
	}
	// write z to disk
	tree.writeNewNode(z)

	midKey := y.items[T-1]

	// remove midKey from y
	y.items = y.items[:T-1]

	// insert z into x at i+1
	x.children = slices.Insert(x.children, i+1, z.offset)
	// insert midKey into x at i
	x.items = slices.Insert(x.items, i, midKey)

	// write y to disk
	tree.writeNodeAtOffset(y)
	// write x to disk
	tree.writeNodeAtOffset(x)
}

func (tree *Btree) InsertNonFull(x *Node, key *Item) {
	idx, _ := slices.BinarySearchFunc(x.items, key, func(a, b *Item) int {
		return a.Compare(b)
	})
	if x.isLeaf {
		// insert key into x
		x.items = slices.Insert(x.items, idx, key)
		tree.writeNodeAtOffset(x)
	} else {
		childNode := tree.readNodeAtOffset(x.children[idx])
		if len(childNode.items) == 2*tree.t-1 {
			tree.SplitChild(x, childNode, idx)
			// important! check if key is in the new right node
			if key.Compare(x.items[idx]) > 0 {
				idx++
			}
		}
		childNode = tree.readNodeAtOffset(x.children[idx])
		tree.InsertNonFull(childNode, key)
	}
}
func (tree *Btree) searchNode(node *Node, key *Item) *Item {
	idx, _ := slices.BinarySearchFunc(node.items, key, func(a, b *Item) int {
		return a.Compare(b)
	})
	if idx < len(node.items) && node.items[idx].Compare(key) == 0 {
		return node.items[idx]
	}
	if node.isLeaf {
		return nil
	}
	childNode := tree.readNodeAtOffset(node.children[idx])
	return tree.searchNode(childNode, key)
}

func (tree *Btree) Search(key *Item) *Item {
	rootNode := tree.readNodeAtOffset(tree.rootOffset)
	return tree.searchNode(rootNode, key)
}
