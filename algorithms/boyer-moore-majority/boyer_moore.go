package boyer_moore_majority

func MajorityElement(arr []int) int {
	count := 0
	majority := -1
	for _, x := range arr {
		if count == 0 {
			majority = x
			count = 1
			continue
		}

		if majority == x {
			count++
		} else {
			count--
		}
	}
	return majority
}
