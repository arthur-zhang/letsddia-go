package hash_ring

import (
	"errors"
	"math"
)

type Key int

type Node struct {
	hashValue int
	prev      *Node
	next      *Node
	resources map[Key]bool
}

func NewNode(hashValue int) *Node {
	return &Node{
		hashValue: hashValue,
		prev:      nil,
		next:      nil,
		resources: make(map[Key]bool),
	}
}

type HashRing struct {
	head *Node
	k    int
	min  int
	max  int
}

func NewHashRing(k int) *HashRing {
	return &HashRing{
		head: nil,
		k:    k,
		min:  0,
		max:  powInt(2, k) - 1,
	}
}

func powInt(base, exp int) int {
	if exp == 0 {
		return 1
	}
	if exp == 1 {
		return base
	}
	return int(math.Pow(float64(base), float64(exp)))
}

// distance
func (hr *HashRing) distance(a, b int) int {
	if a == b {
		return 0
	}
	if a < b {
		return b - a
	}
	return powInt(2, hr.k) - (a - b)
}

func (hr *HashRing) legalRange(hashValue int) bool {
	return hashValue >= hr.min && hashValue <= hr.max
}

// lookupNode find a node that equal or first larger than hashValue
func (hr *HashRing) lookupNode(hashValue int) *Node {
	if !hr.legalRange(hashValue) {
		return nil
	}
	tmp := hr.head
	if tmp == nil {
		return nil
	}
	for hr.distance(tmp.hashValue, hashValue) > hr.distance(tmp.next.hashValue, hashValue) {
		tmp = tmp.next
	}
	if tmp.hashValue == hashValue {
		return tmp
	}
	// always return the first larger node
	return tmp.next
}

func (hr *HashRing) AddNode(hashValue int) {
	if !hr.legalRange(hashValue) {
		return
	}
	newNode := NewNode(hashValue)
	if hr.head == nil {
		hr.head = newNode
		newNode.prev = newNode
		newNode.next = newNode
		println("add node", hashValue, "to empty ring")
		return
	}
	// find the first node that larger than hashValue
	tmp := hr.lookupNode(hashValue)
	newNode.next = tmp
	newNode.prev = tmp.prev
	newNode.prev.next = newNode
	newNode.next.prev = newNode
	println("add node", hashValue, "to ring, prev:", newNode.prev.hashValue, "next:", newNode.next.hashValue)
	hr.moveResources(newNode, newNode.next, false)
	if hashValue < hr.head.hashValue {
		hr.head = newNode
	}
}
func hash(s string) int {
	hashValue := 0
	for _, c := range s {
		hashValue = hashValue*31 + int(c&0xff)
	}
	return hashValue
}
func hashInt(i int) int {
	return i
}
func (hr *HashRing) AddResource(value int) error {
	valueHash := hashInt(value) % hr.max
	node := hr.lookupNode(valueHash)
	if node == nil {
		println("no node found for", value)
		return errors.New("no node found")
	}
	node.resources[Key(value)] = true
	return nil

}
func (hr *HashRing) moveResources(dest, orig *Node, forceMove bool) {

	toDelete := make([]Key, 0)
	for k, _ := range orig.resources {
		valueHash := hashInt(int(k)) % hr.max
		// if the distance between valueHash and dest is smaller than orig, move it
		if forceMove || hr.distance(valueHash, dest.hashValue) < hr.distance(valueHash, orig.hashValue) {
			dest.resources[k] = true
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(orig.resources, k)
	}
}

func (hr *HashRing) DeleteNode(hashValue int) {

	tmp := hr.lookupNode(hashValue)
	if tmp == nil {
		return
	}
	if tmp.hashValue != hashValue {
		println("no node found for", hashValue)
		return
	}
	hr.moveResources(tmp.next, tmp, true)
	tmp.prev.next = tmp.next
	tmp.next.prev = tmp.prev
	if hr.head.hashValue == hashValue {
		hr.head = tmp.next
		if hr.head == hr.head.next {
			hr.head = nil
		}
	}
}
func (hr *HashRing) DebugPrint() {
	println("************************")
	if hr.head == nil {
		println("empty ring")
	}
	tmp := hr.head
	for {
		print("node:", tmp.hashValue, " resources:")
		for k, _ := range tmp.resources {
			print(" ", k)
		}
		println()
		if tmp.next == hr.head {
			break
		}
		tmp = tmp.next
	}
}
