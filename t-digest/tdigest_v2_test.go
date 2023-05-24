package t_digest

import (
	"math"
	"testing"
)

func TestNewTDigestV2(t *testing.T) {
	td := NewTDigestV2(5)
	if !math.IsNaN(td.Quantile(0.5)) {
		t.Errorf("Expected NaN for empty TDigest")
	}
	td.Insert(1.0)
	if td.Quantile(0.5) != 1.0 {
		t.Errorf("Expected 1.0 for TDigest with one element")
	}

	for i := 2.0; i <= 10; i++ {
		td.Insert(i)
	}
	if td.Quantile(0.5) < 4.0 || td.Quantile(0.5) > 6.0 {
		t.Errorf("Expected value around 5.0 for TDigest with ten elements")
	}
	td = NewTDigestV2(1000)
	for i := 1.0; i <= 1000; i++ {
		td.Insert(i)
	}
	if td.Quantile(0.5) < 490.0 || td.Quantile(0.5) > 510.0 {
		t.Errorf("Expected value around 5000.0 for TDigest with ten thousand elements")
	}
}
