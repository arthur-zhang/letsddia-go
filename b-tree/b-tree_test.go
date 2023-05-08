package b_tree

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
}
func TestInsert(t *testing.T) {
	tree := NewBtree(2)
	arr := []int{'F', 'S', 'Q', 'K', 'C', 'L', 'H', 'T', 'V', 'W', 'M', 'R', 'N', 'P', 'A', 'B', 'X', 'Y', 'D', 'Z', 'E'}

	for _, v := range arr {
		fmt.Printf("%c\n", v)
		tree.Insert(v)
		fmt.Printf("%c done\n", v)
	}
}

func TestNodeSpit(t *testing.T) {

	y := NewNode(true)
	for i := 0; i < 7; i++ {
		y.children = nil
		y.keys = append(y.keys, i)
	}
	x := NewNode(false)
	x.keys = append(x.keys, -1)
	x.keys = append(x.keys, 8)
	//x.children = append(x.children, c0)
	x.children = append(x.children, nil)
	x.children = append(x.children, y)
	x.children = append(x.children, nil)
	x.splitChild(y, 1)
	println(x.keys[0])

}
