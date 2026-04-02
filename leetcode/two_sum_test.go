package main

import "testing"

func TestTwoSum(t *testing.T) {
	nums := []int{2, 7, 11, 15}
	target := 9
	got := twoSum(nums, target)
	want := []int{0, 1}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}
