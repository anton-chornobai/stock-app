package tests

import "testing"

func FindMax(arr []int) int {
	
	max := 0;

	for _, val := range arr {
		if val > max {
			max = val
		}
	}

	return max;
}

func TestFindMax(t *testing.T) {
	arr := []int{1, 2, 3, -5, -2, 3, 10, 15, 9, 2}
	max := FindMax(arr)

	if max != 15 {
		t.Errorf("expected %d", 15)
	}
}