package bitset

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitSet_Set(t *testing.T) {
	bs := New(128)
	bs.Set(1)
	bs.Set(10)
	bs.Set(50)
	bs.Set(2)
	assert.True(t, bs.IsSet(1))
	assert.True(t, bs.IsSet(10))
	assert.True(t, bs.IsSet(50))
	assert.True(t, bs.IsSet(2))
}
func TestBitSet_Set2(t *testing.T) {
	bs := New(1 << 32)
	assert.Equal(t, bs.Len(), uint64(1<<32))
	assert.False(t, bs.IsSet(0))
	bs.Set(1)
	assert.True(t, bs.IsSet(1))
	bs.Set(1<<32 - 1)
	assert.True(t, bs.IsSet(1<<32-1))
}
func TestBitSet_Clear(t *testing.T) {
	bs := New(128)
	bs.Set(1)
	bs.Set(50)
	assert.True(t, bs.IsSet(1))
	bs.Clear(1)
	assert.False(t, bs.IsSet(1))
	assert.True(t, bs.IsSet(50))
	bs.Clear(50)
	assert.False(t, bs.IsSet(50))
}
