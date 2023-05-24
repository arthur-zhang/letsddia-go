package lo

import "testing"

type Foo struct {
	x int
}

func TestSumBy(t *testing.T) {
	arr := []Foo{{1}, {2}, {3}}
	sum := SumBy(arr, func(f Foo) int {
		return f.x
	})
	println(sum)
}
