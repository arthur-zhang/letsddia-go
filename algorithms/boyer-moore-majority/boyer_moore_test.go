package boyer_moore_majority

import (
	"testing"
)

func TestMajorityElement(t *testing.T) {
	// Test case 1: Majority element is in the middle
	nums1 := []int{2, 4, 7, 4, 2, 4, 4}
	expected1 := 4
	result1 := MajorityElement(nums1)
	if result1 != expected1 {
		t.Errorf("Expected %d, but got %d", expected1, result1)
	}

	// Test case 2: Majority element is at the beginning
	nums2 := []int{3, 3, 3, 2, 2, 4, 4, 3, 3}
	expected2 := 3
	result2 := MajorityElement(nums2)
	if result2 != expected2 {
		t.Errorf("Expected %d, but got %d", expected2, result2)
	}

	// Test case 3: Majority element is at the end
	nums3 := []int{0, 1, 2, 1, 3, 1, 4, 1, 1, 1}
	expected3 := 1
	result3 := MajorityElement(nums3)
	if result3 != expected3 {
		t.Errorf("Expected %d, but got %d", expected3, result3)
	}

	// Test case 4: Majority element is the only element
	nums4 := []int{5}
	expected4 := 5
	result4 := MajorityElement(nums4)
	if result4 != expected4 {
		t.Errorf("Expected %d, but got %d", expected4, result4)
	}

	// Test case 5: Empty array
	var nums5 []int
	expected5 := -1 // No majority element, return -1 as a sentinel value
	result5 := MajorityElement(nums5)
	if result5 != expected5 {
		t.Errorf("Expected %d, but got %d", expected5, result5)
	}
}
