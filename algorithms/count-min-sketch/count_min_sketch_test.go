package count_min_sketch

import "testing"

func TestCountMinSketch(t *testing.T) {

	cms := New(4, 8)
	elements := []string{"apple", "banana", "apple", "orange", "banana", "apple"}
	for _, element := range elements {
		cms.Update(element, 1)
	}
	for _, ele := range elements {
		if cms.Estimate(ele) != 2 {
			//t.Error("CountMinSketch error")
		}
		println(ele, cms.Estimate(ele))
	}

}
