package t_digest

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestNewTDigestV2(t *testing.T) {
	td := NewTDigestV2()
	if !math.IsNaN(td.Quantile(0.5)) {
		t.Errorf("Expected NaN for empty TDigest")
	}
	td.Insert(1.0)
	//if td.Quantile(0.5) != 1.0 {
	//	t.Errorf("Expected 1.0 for TDigest with one element")
	//}

	for i := 2.0; i <= 10; i++ {
		td.Insert(i)
	}
	fmt.Printf("%f\n", td.Quantile(0.5))
	if td.Quantile(0.5) < 4.0 || td.Quantile(0.5) > 6.0 {
		t.Errorf("Expected value around 5.0 for TDigest with ten elements %f", td.Quantile(0.5))
	}
	td = NewTDigestV2()
	for i := 1.0; i <= 10000; i++ {
		td.Insert(i)
	}
	fmt.Printf("%f\n", td.Quantile(0.5))
	if td.Quantile(0.5) < 4900.0 || td.Quantile(0.5) > 5100.0 {
		t.Errorf("Expected value around 5000.0 for TDigest with ten thousand elements %f", td.Quantile(0.5))
	}
}

func Test_getPotential(t *testing.T) {
	for i := 0.0; i <= 1; i += 0.1 {
		a := getK(i)
		fmt.Printf("%f: %f\n", i, a)
	}
}
func TestRandom(t *testing.T) {
	td := NewTDigestV2()
	for i := 1; i <= 10000; i++ {
		td.Insert(float64(rand.Uint32() % 1000))
	}
	fmt.Printf("%f\n", td.Quantile(0.9))
}
