package roaring_bitmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayContainer_Add(t *testing.T) {
	c := ArrayContainer{}
	c.Add(1)
	c.Add(4)
	c.Add(2)
	c.Add(3)
	for _, item := range c.content {
		println(item)
	}
	assert.True(t, c.Contains(1))
	assert.True(t, c.Contains(2))
	assert.True(t, c.Contains(3))
	assert.True(t, c.Contains(4))
	println("----")
	c.Remove(1)
	for _, item := range c.content {
		println(item)
	}
	assert.False(t, c.Contains(1))
	assert.True(t, c.Contains(2))
	assert.True(t, c.Contains(3))
	assert.True(t, c.Contains(4))
}

func TestArrayToBitmap(t *testing.T) {
	var c Container
	c = &ArrayContainer{}
	i := 0
	for i = 0; i < ArrayToBitmapThreshold; i++ {
		c, _ = c.Add(uint16(i))
	}
	for i = 0; i < ArrayToBitmapThreshold; i++ {
		assert.True(t, c.Contains(uint16(i)))
	}
	for ; i < ArrayToBitmapThreshold*2; i++ {
		assert.False(t, c.Contains(uint16(i)))
	}
	assert.IsType(t, &ArrayContainer{}, c)

	// add new one
	c, _ = c.Add(uint16(i))
	assert.IsType(t, &BitmapContainer{}, c)

}

func TestRunContainer_Add(t *testing.T) {
	var c Container
	c = NewRunContainer()
	c, _ = c.Add(11)
	c, _ = c.Add(12)
	c, _ = c.Add(13)
	c, _ = c.Add(14)
	c, _ = c.Add(9)
	c, _ = c.Add(20)
	for i, run := range c.(RunContainer) {
		println(i, run.start, run.end)
	}
}
