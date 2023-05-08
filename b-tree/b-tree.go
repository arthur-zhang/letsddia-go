package b_tree

const T = 2

type Node struct {
	IsLeaf   bool
	keys     []int
	children []*Node
}

func NewNode(isLeaf bool) *Node {
	return &Node{
		IsLeaf:   isLeaf,
		keys:     make([]int, 0, T*2-1),
		children: make([]*Node, 0, T*2),
	}
}

type Btree struct {
	root *Node
}

func NewBtree(t int) Btree {
	return Btree{
		root: NewNode(true),
	}
}

func (b *Btree) search(x *Node, key int) (*Node, int) {
	if x == nil {
		return nil, 0
	}

	i := 0
	// Search for the key in the current node and find the first index i
	// such that key <= x.keys[i]
	for i <= len(x.keys) && key > x.keys[i] {
		i++
	}
	// Check if the key is found at index i
	if i < len(x.keys) && key == x.keys[i] {
		return x, i
	}
	// If the current node is a leaf and the key is not found,
	// return nil, 0
	if x.IsLeaf {
		return nil, 0
	}
	// Recursively search the children nodes
	return b.search(x.children[i], key)
}

func (b *Btree) Insert(key int) {
	// If the root is full, split it
	if len(b.root.keys) == 2*T-1 {
		// New root
		s := NewNode(false)
		// Make the original root the child of the new root
		s.children = append(s.children, b.root)
		// Split the original root
		s.splitChild(b.root, 0)
		// Insert the key
		s.insertNonFull(key)
		b.root = s
	} else {
		b.root.insertNonFull(key)
	}
}

func (x *Node) insertNonFull(key int) {
	i := len(x.keys) - 1
	if x.IsLeaf {
		// Traverse from the end of keys and find the first index i such that
		// key >= x.keys[i]
		x.keys = append(x.keys, 0)
		for i >= 0 && key < x.keys[i] {
			// Move the values in keys right by one spot
			x.keys[i+1] = x.keys[i]
			i--
		}
		x.keys[i+1] = key
	} else {
		// Traverse from the end of keys and find the first index i such that
		// key <= x.keys[i]
		for i >= 0 && key < x.keys[i] {
			i--
		}
		i++
		// If the child node is full, split it
		if len(x.children[i].keys) == 2*T-1 {
			x.splitChild(x.children[i], i)
			if key > x.keys[i] {
				i++
			}
		}
		x.children[i].insertNonFull(key)
	}
}

func (x *Node) splitChild(y *Node, i int) {
	// The degree of the new node is the half of the maximum capacity of keys of
	// full node y
	t := T
	// The degree of the new node z is the half of the maximum capacity of keys
	// of full node y
	z := NewNode(y.IsLeaf)

	// Copy the second half of keys from y to z
	z.keys = append(z.keys, y.keys[t:]...)
	// If y is not a leaf node, copy the second half of children from y to z and
	// delete them from y
	if !z.IsLeaf {
		z.children = append(z.children, y.children[t:]...)
		y.children = y.children[:t]
	}

	midKey := y.keys[t-1]
	// Remove the middle key from y
	y.keys = y.keys[:t-1]
	// Add a nil child to x to create room
	x.children = append(x.children, nil)
	// Shift x.children[i+1:] right one spot
	copy(x.children[i+1:], x.children[i:])
	// Insert new node z at x.children[i+1]
	x.children[i+1] = z
	// Add a zero key to x to create room
	x.keys = append(x.keys, 0)
	// Shift x.keys[i:] right one spot
	copy(x.keys[i+1:], x.keys[i:])
	// Put the middle key of y into x.keys[i]
	x.keys[i] = midKey
}
