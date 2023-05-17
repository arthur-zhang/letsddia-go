package hash_ring

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistance(t *testing.T) {
	hr := NewHashRing(3)
	assert.Equal(t, 2, hr.distance(6, 0))
	assert.Equal(t, 6, hr.distance(0, 6))
	assert.Equal(t, 3, hr.distance(3, 6))
}
func TestHashRing_AddNode(t *testing.T) {
	hr := NewHashRing(5)
	hr.AddNode(12)
	hr.AddNode(18)
	hr.AddNode(5)
}

func TestHashRing(t *testing.T) {
	hr := NewHashRing(5)
	hr.AddNode(12)
	hr.AddNode(18)
	hr.AddResource(24)
	hr.AddResource(21)
	hr.AddResource(16)
	hr.AddResource(23)
	hr.AddResource(2)
	hr.AddResource(29)
	hr.AddResource(28)
	hr.AddResource(7)
	hr.AddResource(10)
	hr.DebugPrint()
	hr.AddNode(5)
	hr.AddNode(27)
	hr.AddNode(30)
	hr.DebugPrint()
	hr.DeleteNode(12)
	hr.DebugPrint()
}
