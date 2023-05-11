package b_tree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	tree := NewBtree()
	arr := []int{'F', 'S', 'Q', 'K', 'C', 'L', 'H', 'T', 'V', 'W', 'M', 'R', 'N', 'P', 'A', 'B', 'X', 'Y', 'D', 'Z', 'E'}

	for _, v := range arr {
		tree.Insert(v)
	}
	node, idx := tree.Search('Y')
	assert.True(t, node != nil)
	assert.Equal(t, 1, idx)

	node, idx = tree.Search('G')
	assert.True(t, node == nil)
}
func TestInsert(t *testing.T) {
	tree := NewBtree()
	arr := []int{'F', 'S', 'Q', 'K', 'C', 'L', 'H', 'T', 'V', 'W', 'M', 'R', 'N', 'P', 'A', 'B', 'X', 'Y', 'D', 'Z', 'E'}

	for _, v := range arr {
		tree.Insert(v)
	}
	r := tree.root
	assert.EqualValues(t, []int{'K', 'Q'}, r.keys)
	l := len(r.children)
	assert.Equal(t, 3, l)
	assert.EqualValues(t, []int{'B', 'F'}, r.children[0].keys)
	assert.EqualValues(t, []int{'M'}, r.children[1].keys)
	assert.EqualValues(t, []int{'T', 'W'}, r.children[2].keys)

	assert.Equal(t, 3, len(r.children[0].children))
	assert.EqualValues(t, []int{'A'}, r.children[0].children[0].keys)
	assert.EqualValues(t, []int{'C', 'D', 'E'}, r.children[0].children[1].keys)
	assert.EqualValues(t, []int{'H'}, r.children[0].children[2].keys)

	assert.Equal(t, 2, len(r.children[1].children))
	assert.EqualValues(t, []int{'L'}, r.children[1].children[0].keys)
	assert.EqualValues(t, []int{'N', 'P'}, r.children[1].children[1].keys)

	assert.Equal(t, 3, len(r.children[2].children))
	assert.EqualValues(t, []int{'R', 'S'}, r.children[2].children[0].keys)
	assert.EqualValues(t, []int{'V'}, r.children[2].children[1].keys)
	assert.EqualValues(t, []int{'X', 'Y', 'Z'}, r.children[2].children[2].keys)
}
