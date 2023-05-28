package lo

import "golang.org/x/exp/constraints"

func SumBy[T any, R constraints.Float | constraints.Integer | constraints.Complex](arr []T, f func(T) R) R {
	var sum R = 0
	for _, x := range arr {
		sum += f(x)
	}
	return sum
}
func FindIndexOf[T any](arr []T, f func(T) bool) (T, int, bool) {
	for i, item := range arr {
		if f(item) {
			return item, i, true
		}
	}
	var result T
	return result, -1, false
}

func Map[T any, R any](arr []T, f func(T, int) R) []R {
	result := make([]R, len(arr))
	for i, item := range arr {
		result[i] = f(item, i)
	}
	return result
}
