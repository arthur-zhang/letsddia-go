package roaring_bitmap

import (
	"bitset"
	"sort"
)

type Container interface {
	Add(x uint16) (Container, bool)
	Contains(x uint16) bool
	Remove(x uint16) (Container, bool)
	Len() int
}

// ArrayContainer ----------------------------------
type ArrayContainer struct {
	content []uint16
}

const ArrayToBitmapThreshold = 4096
const ArrayToRunThreshold = 64

func NewArrayContainer() *ArrayContainer {
	return &ArrayContainer{
		content: make([]uint16, 0),
	}
}

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

	if len(c.content) > ArrayToBitmapThreshold {
		return c.convertToBitmapContainer(), true
	}
	if c.hasEnoughRuns() {
		return c.convertToRunContainer(), true
	}
	return c, true
}

func (c *ArrayContainer) convertToRunContainer() Container {
	rc := NewRunContainer()
	for i := 0; i < len(c.content); i++ {
		rc.Add(c.content[i])
	}
	return &rc
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

func (c *ArrayContainer) Len() int {
	return len(c.content)
}

func (c *ArrayContainer) hasEnoughRuns() bool {
	count := 1
	for i := 1; i < len(c.content); i++ {
		if c.content[i]-c.content[i-1] == 1 {
			count++
			if count > ArrayToRunThreshold {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
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
func (b *BitmapContainer) Add(x uint16) (Container, bool) {
	b.bitmap.Set(uint64(x))
	return b, true
}

func (b *BitmapContainer) Contains(x uint16) bool {
	return b.bitmap.IsSet(uint64(x))
}

func (b *BitmapContainer) Remove(x uint16) (Container, bool) {
	_ = b.bitmap.Clear(uint64(x))
	return b, true
}

func (b *BitmapContainer) Len() int {
	return b.bitmap.Count()
}

// RunContainer ----------------------------------
type RunContainer struct {
	runs []Run
}

type Run struct {
	start uint16
	end   uint16
}

func NewRunContainer() RunContainer {
	return RunContainer{
		runs: make([]Run, 0),
	}
}

func (rc *RunContainer) Add(x uint16) (Container, bool) {
	n := len(rc.runs)
	if n == 0 {
		rc.runs = append(rc.runs, Run{x, x})
		return rc, true
	}
	// find the first run that starts after x
	i := sort.Search(len(rc.runs), func(i int) bool {
		return rc.runs[i].start > x
	})

	// x is already in the container
	if i > 0 && rc.runs[i-1].start <= x && x <= rc.runs[i-1].end {
		return rc, false
	}
	// can add to previous run
	if i > 0 && rc.runs[i-1].end == x-1 {
		rc.runs[i-1].end = x
		// merge with next
		if i < n && rc.runs[i].start == x+1 {
			rc.runs[i-1].end = rc.runs[i].end
			rc.runs = append(rc.runs[:i], rc.runs[i+1:]...)
			return rc, true
		}
		return rc, true
	}

	// 8,9  11,12,13
	// can add to next run
	if i < n && rc.runs[i].start == x+1 {
		rc.runs[i].start = x
		// merge with previous
		if i > 0 && rc.runs[i-1].end == x-1 {
			rc.runs[i-1].end = rc.runs[i].end
			rc.runs = append(rc.runs[:i], rc.runs[i+1:]...)
			return rc, true
		}
		return rc, true
	}

	newRun := Run{x, x}
	rc.runs = append(rc.runs, newRun)
	if i < len(rc.runs)-1 {
		copy(rc.runs[i+1:], rc.runs[i:])
	}
	rc.runs[i] = newRun
	return rc, true
}

func (rc *RunContainer) Contains(x uint16) bool {
	idx := sort.Search(len(rc.runs), func(i int) bool {
		return rc.runs[i].end >= x
	})
	return idx < len(rc.runs) && rc.runs[idx].start <= x && x <= rc.runs[idx].end
}

func (rc *RunContainer) Remove(x uint16) (Container, bool) {
	idx := sort.Search(len(rc.runs), func(i int) bool {
		return rc.runs[i].start > x
	})
	idx -= 1
	if idx < 0 || idx >= len(rc.runs) {
		return rc, false
	}
	if rc.runs[idx].start > x || rc.runs[idx].end < x {
		return rc, false
	}
	if rc.runs[idx].start == x {
		if rc.runs[idx].end == x {
			rc.runs = append(rc.runs[:idx], rc.runs[idx+1:]...)
			return rc, true
		}
		rc.runs[idx].start += 1
		return rc, true
	}
	if rc.runs[idx].end == x {
		rc.runs[idx].end -= 1
		return rc, true
	}

	// now we should split
	newRun := Run{x + 1, rc.runs[idx].end}
	rc.runs[idx].end = x - 1
	rc.runs = append(rc.runs, newRun)
	if idx < len(rc.runs)-1 {
		copy(rc.runs[idx+2:], rc.runs[idx+1:])
	}
	rc.runs[idx+1] = newRun
	return rc, true
}

func (rc *RunContainer) Len() int {
	count := 0
	for _, r := range rc.runs {
		count += int(r.end - r.start + 1)
	}
	return count
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
		c = NewArrayContainer()
	}
	c, _ = c.Add(low)
	rb.highLowContainer[high] = c
}

// Remove todo optimize container type when remove
func (rb *RoaringBitmap) Remove(x uint32) {
	high := uint16(x >> 16)
	low := uint16(x & 0xFFFF)
	c, ok := rb.highLowContainer[high]
	if !ok {
		return
	}
	c, _ = c.Remove(low)
	rb.highLowContainer[high] = c
}
func (rb *RoaringBitmap) Len() int {
	count := 0
	for _, c := range rb.highLowContainer {
		count += c.Len()
	}
	return count
}
