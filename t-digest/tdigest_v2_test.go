package t_digest

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNewTDigestV2(t_ *testing.T) {
	t := NewTDigestV2(100)
	for i := 1000; i >= 1; i-- {
		t.Insert(float64(rand.Uint32() % 100))
	}
	fmt.Printf("%f\n", t.Quantile(0.4))
}
