package roaring_bitmap

import (
	"bitset"
	"sort"
)

type Container interface {
	Add(x uint16) (Container, bool)
	Contains(x uint16) bool
	Remove(x uint16) (Container, bool)
}

// ArrayContainer ----------------------------------
type ArrayContainer struct {
	content []uint16
}

const ArrayToBitmapThreshold = 4096

func (c *ArrayContainer) Add(x uint16) (Container, bool) {
	idx := sort.Search(len(c.content), func(i int) bool {
		return c.content[i] >= x
	})
	if idx < len(c.content) && c.content[idx] == x {
		return c, false
	}

	c.content = append(c.content, x)
	copy(c.content[idx+1:], c.content[idx:])
	c.content[idx] = x

	if len(c.content) <= ArrayToBitmapThreshold {
		return c, true
	}
	return c.convertToBitmapContainer(), true
}

func (c *ArrayContainer) convertToBitmapContainer() Container {
	newC := NewBitmapContainer()
	for _, item := range c.content {
		newC.Add(item)
	}
	return newC
}

func (c *ArrayContainer) Contains(x uint16) bool {
	idx := sort.Search(len(c.content), func(i int) bool {
		return c.content[i] >= x
	})
	return idx < len(c.content) && c.content[idx] == x
}

func (c *ArrayContainer) Remove(x uint16) (Container, bool) {
	idx := sort.Search(len(c.content), func(i int) bool {
		return c.content[i] >= x
	})
	if idx < len(c.content) && c.content[idx] == x {
		c.content = append(c.content[:idx], c.content[idx+1:]...)
		return c, true
	}
	return c, false
}

// BitmapContainer ----------------------------------
type BitmapContainer struct {
	bitmap bitset.BitSet
}

func NewBitmapContainer() Container {
	return &BitmapContainer{
		bitmap: bitset.New(65536),
	}
}
func (b BitmapContainer) Add(x uint16) (Container, bool) {
	b.bitmap.Set(uint64(x))
	return b, true
}

func (b BitmapContainer) Contains(x uint16) bool {
	return b.bitmap.IsSet(uint64(x))
}

func (b BitmapContainer) Remove(x uint16) (Container, bool) {
	panic("implement me")
}

// RunContainer ----------------------------------
type RunContainer []Run

type Run struct {
	start uint16
	end   uint16
}

func NewRunContainer() RunContainer {
	return make(RunContainer, 0)
}

func (rc RunContainer) Add(x uint16) (Container, bool) {

	n := len(rc)
	if n == 0 {
		return append(rc, Run{x, x}), true
	}
	// find the first run that starts after x
	i := sort.Search(len(rc), func(i int) bool {
		return rc[i].start > x
	})

	// x is already in the container
	if i > 0 && rc[i-1].start <= x && x <= rc[i-1].end {
		return rc, false
	}
	// can add to previous run
	if i > 0 && rc[i-1].end == x-1 {
		rc[i-1].end = x
		// merge with next
		if i < n && rc[i].start == x+1 {
			rc[i-1].end = rc[i].end
			return append(rc[:i], rc[i+1:]...), true
		}
		return rc, true
	}

	// 8,9  11,12,13
	// can add to next run
	if i < n && rc[i].start == x+1 {
		rc[i].start = x
		// merge with previous
		if i > 0 && rc[i-1].end == x-1 {
			rc[i-1].end = rc[i].end
			return append(rc[:i], rc[i+1:]...), true
		}
		return rc, true
	}

	newRun := Run{x, x}
	rc = append(rc, newRun)
	if i < len(rc)-1 {
		copy(rc[i+1:], rc[i:])
	}
	rc[i] = newRun
	return rc, true
}

func (rc RunContainer) Contains(x uint16) bool {
	idx := sort.Search(len(rc), func(i int) bool {
		return rc[i].end >= x
	})
	return idx < len(rc) && rc[idx].start <= x && x <= rc[idx].end
}

func (rc RunContainer) Remove(x uint16) (Container, bool) {
	panic("implement me")
}

// RoaringBitmap ----------------------------------
type RoaringBitmap struct {
	//highLowContainer []Entry
	highLowContainer map[uint16]Container
}
type Entry struct {
	key   uint16
	value Container
}

func New() *RoaringBitmap {
	return &RoaringBitmap{
		highLowContainer: make(map[uint16]Container),
	}
}

func (rb *RoaringBitmap) Add(x uint32) {
	high := uint16(x >> 16)
	low := uint16(x & 0xFFFF)
	c, ok := rb.highLowContainer[high]
	if !ok {
		c = &ArrayContainer{}
	}
	c, _ = c.Add(low)
	rb.highLowContainer[high] = c
}
