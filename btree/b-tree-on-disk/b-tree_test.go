package b_tree_on_disk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"testing"
)

func TestPageToNode(t *testing.T) {
	data := []byte{
		0x01,                                           // isLeaf
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // Number of keys.
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, // 4096  (2nd Page)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, // 8192  (3rd Page)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, // 12288 (4th Page)
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, // "hello"
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, // "world"
		0x66, 0x6f, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // "foo"
		0x62, 0x61, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // "bar"
	}
	pageData := make([]byte, PageSize)
	copy(pageData, data)
	page := &Page{
		data: [PageSize]byte(pageData),
	}
	node := page.ToNode()
	assert.True(t, node.isLeaf)
	assert.Equal(t, 2, len(node.items))
	assert.Equal(t, 3, len(node.children))
	assert.Equal(t, uint64(4096), node.children[0])
	assert.Equal(t, uint64(8192), node.children[1])
	assert.Equal(t, uint64(12288), node.children[2])
	assert.Equal(t, "hello", string(node.items[0].key[:5]))
	assert.Equal(t, "world", string(node.items[0].value[:5]))
	assert.Equal(t, "foo", string(node.items[1].key[:3]))
	assert.Equal(t, "bar", string(node.items[1].value[:3]))
}

func TestNode2Page(t *testing.T) {
	data := []byte{
		0x01,                                           // isLeaf
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // Number of keys.
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, // 4096  (2nd Page)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, // 8192  (3rd Page)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x00, // 12288 (4th Page)
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, // "hello"
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, // "world"
		0x66, 0x6f, 0x6f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // "foo"
		0x62, 0x61, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // "bar"
	}
	pageData := make([]byte, PageSize)
	copy(pageData, data)

	node := &Node{
		isLeaf: true,
		items: []*Item{
			{
				key:   [N]byte{'h', 'e', 'l', 'l', 'o'},
				value: [N]byte{'w', 'o', 'r', 'l', 'd'},
			},
			{
				key:   [N]byte{'f', 'o', 'o', 0, 0, 0, 0, 0, 0, 0},
				value: [N]byte{'b', 'a', 'r'},
			},
		},
		children: []int64{4096, 8192, 12288},
	}
	page := node.ToPage()
	assert.Equal(t, pageData, page.data[:])
}
func (tree *Btree) Traverse() {
	rootNode := tree.readNodeAtOffset(tree.rootOffset)
	tree.inOrderVisitNode(rootNode)
}
func (tree *Btree) inOrderVisitNode(node *Node) {
	if node == nil {
		return
	}

	for i := 0; i < len(node.items); i++ {
		if !node.isLeaf {
			childNode := tree.readNodeAtOffset(node.children[i])
			tree.inOrderVisitNode(childNode)
		}
		fmt.Printf("%s ", decodeString(node.items[i].key))
	}

	if !node.isLeaf {
		childNode := tree.readNodeAtOffset(node.children[len(node.children)-1])
		tree.inOrderVisitNode(childNode)
	}
}

func decodeString(data [N]byte) string {
	idx := slices.Index(data[:], 0)
	if idx == -1 {
		return string(data[:])
	} else {
		return string(data[:idx])
	}
}
func TestInsert(t *testing.T) {
	btree := NewBtree("/tmp/btree.db", 2)
	arr := []byte{'F', 'S', 'Q', 'K', 'C', 'L', 'H', 'T', 'V', 'W', 'M', 'R', 'N', 'P', 'A', 'B', 'X', 'Y', 'D', 'Z', 'E'}
	for _, v := range arr {
		item := &Item{
			key:   [N]byte{v},
			value: [N]byte{v, byte('0'), byte('1')},
		}
		btree.Insert(item)
	}
	btree.Traverse()
	key := Item{key: [N]byte{'A'}}
	item := btree.Search(&key)
	assert.True(t, item != nil)
	assert.Equal(t, "A", decodeString(item.key))
	assert.Equal(t, "A01", decodeString(item.value))

	key = Item{key: [N]byte{'a'}}
	item = btree.Search(&key)
	assert.True(t, item == nil)
}
