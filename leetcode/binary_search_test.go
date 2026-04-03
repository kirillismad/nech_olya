package leetcode

import "testing"

func TestBinarySearch(t *testing.T) {
	t.Run("tc 1", func(t *testing.T) {
		nums := []int{-1, 0, 3, 5, 9, 12}
		target := 9
		got := search(nums, target)
		want := 4
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("tc 2", func(t *testing.T) {
		nums := []int{-1, 0, 3, 5, 9, 12}
		target := 2
		got := search(nums, target)
		want := -1
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
