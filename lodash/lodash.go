package lo

import "golang.org/x/exp/constraints"

func SumBy[T any, R constraints.Float | constraints.Integer | constraints.Complex](arr []T, f func(T) R) R {
	var sum R = 0
	for _, x := range arr {
		sum += f(x)
	}
	return sum
}
