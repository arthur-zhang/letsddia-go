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
	rc := NewRunContainer()
	c = &rc
	c, _ = c.Add(11)

	assert.Equal(t, (c.(*RunContainer)).runs[0].start, uint16(11))
	assert.Equal(t, (c.(*RunContainer)).runs[0].end, uint16(11))
	c, _ = c.Add(12)
	assert.Equal(t, len(c.(*RunContainer).runs), 1)
	assert.Equal(t, (c.(*RunContainer)).runs[0].start, uint16(11))
	assert.Equal(t, (c.(*RunContainer)).runs[0].end, uint16(12))
	c, _ = c.Add(13)
	c, _ = c.Add(14)
	assert.Equal(t, len(c.(*RunContainer).runs), 1)
	assert.Equal(t, (c.(*RunContainer)).runs[0].start, uint16(11))
	assert.Equal(t, (c.(*RunContainer)).runs[0].end, uint16(14))

	c, _ = c.Add(9)
	assert.Equal(t, len(c.(*RunContainer).runs), 2)
	assert.Equal(t, (c.(*RunContainer)).runs[0].start, uint16(9))
	assert.Equal(t, (c.(*RunContainer)).runs[0].end, uint16(9))
	assert.Equal(t, (c.(*RunContainer)).runs[1].start, uint16(11))
	assert.Equal(t, (c.(*RunContainer)).runs[1].end, uint16(14))
	c, _ = c.Add(20)
	assert.Equal(t, len(c.(*RunContainer).runs), 3)
	assert.Equal(t, (c.(*RunContainer)).runs[0].start, uint16(9))
	assert.Equal(t, (c.(*RunContainer)).runs[0].end, uint16(9))
	assert.Equal(t, (c.(*RunContainer)).runs[1].start, uint16(11))
	assert.Equal(t, (c.(*RunContainer)).runs[1].end, uint16(14))
	assert.Equal(t, (c.(*RunContainer)).runs[2].start, uint16(20))
	assert.Equal(t, (c.(*RunContainer)).runs[2].end, uint16(20))
	for i, run := range c.(*RunContainer).runs {
		println(i, run.start, run.end)
	}
}

func TestArray2Run(t *testing.T) {
	var c Container
	c = NewArrayContainer()
	for i := 0; i < ArrayToRunThreshold; i++ {
		c, _ = c.Add(uint16(i))
	}
	assert.IsType(t, &ArrayContainer{}, c)
	c, _ = c.Add(uint16(ArrayToRunThreshold))
	assert.IsType(t, &RunContainer{}, c)

	rc := c.(*RunContainer)
	assert.Equal(t, len(rc.runs), 1)
	assert.Equal(t, rc.runs[0].start, uint16(0))
	assert.Equal(t, rc.runs[0].end, uint16(ArrayToRunThreshold))
}

func TestRoaringBitmap_Add(t *testing.T) {
	rb := New()
	for i := 0; i < 100; i++ {
		rb.Add(uint32(i))
	}
	c := rb.highLowContainer[0]
	assert.IsType(t, &RunContainer{}, c)
	assert.Equal(t, rb.Len(), 100)
	for i := 0; i < 100; i++ {
		assert.True(t, c.Contains(uint16(i)))
	}
	assert.False(t, c.Contains(100))

	for i := 0; i < 100; i++ {
		rb.Remove(uint32(i))
	}
	assert.Equal(t, rb.Len(), 0)
}

func TestRoaringBitmap_Add2(t *testing.T) {
	rb := New()
	for i := 0; i < 200; i++ {
		if i%2 == 0 {
			rb.Add(uint32(i))
		}
	}
	assert.Equal(t, rb.Len(), 100)

	c := rb.highLowContainer[0]
	assert.IsType(t, &ArrayContainer{}, c)
	for i := 0; i < 200; i++ {
		if i%2 == 0 {
			assert.True(t, c.Contains(uint16(i)))
		}
	}
	assert.False(t, c.Contains(99))
	assert.False(t, c.Contains(200))

	for i := 0; i < 200; i++ {
		rb.Remove(uint32(i))
	}
	assert.Equal(t, rb.Len(), 0)
}

func TestRoaringBitmap_Add3(t *testing.T) {
	rb := New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			rb.Add(uint32(i))
		}
	}
	assert.Equal(t, rb.Len(), 5000)

	c := rb.highLowContainer[0]
	assert.IsType(t, &BitmapContainer{}, c)
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			assert.True(t, c.Contains(uint16(i)))
		}
	}
	assert.False(t, c.Contains(99))
	assert.False(t, c.Contains(20000))

	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			rb.Remove(uint32(i))
		}
	}
	assert.Equal(t, rb.Len(), 0)
}

func TestRunContainer_Remove(t *testing.T) {
	var c Container
	rc := NewRunContainer()
	c = &rc
	c, _ = c.Add(11)
	assert.Equal(t, c.Len(), 1)
	c, _ = c.Remove(11)
	assert.Equal(t, c.Len(), 0)

	c, _ = c.Add(9)

	c, _ = c.Add(11)
	c, _ = c.Add(12)
	c, _ = c.Add(13)
	c, _ = c.Add(14)
	c, _ = c.Add(20)

	runs := c.(*RunContainer).runs
	assert.Equal(t, len(runs), 3)
	runs = c.(*RunContainer).runs
	c.Remove(0)
	assert.Equal(t, len(runs), 3)
	c.Remove(9)
	runs = c.(*RunContainer).runs
	assert.Equal(t, len(runs), 2)
	assert.Equal(t, runs[0].start, uint16(11))

	c.Remove(13)
	runs = c.(*RunContainer).runs
	assert.Equal(t, len(runs), 3)
	assert.Equal(t, runs[0].start, uint16(11))
	assert.Equal(t, runs[1].start, uint16(14))
	assert.Equal(t, runs[2].start, uint16(20))
}
