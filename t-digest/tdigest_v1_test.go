package t_digest

import (
	"fmt"
	"testing"
)

func TestTDigestV1_Quantile(t1 *testing.T) {
	t := NewTDigestV1()
	for i := 10; i >= 1; i-- {
		t.Insert(float64(i))
	}
	q := t.Quantile(0.8)
	fmt.Printf("%f\n", q)

}
