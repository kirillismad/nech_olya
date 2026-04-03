package leetcode

import "testing"

func TestContainsDuplicate(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		nums := []int{1, 2, 3, 1}
		got := containsDuplicate(nums)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		nums := []int{1, 2, 3, 4}
		got := containsDuplicate(nums)
		want := false
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("tc 3:", func(t *testing.T) {
		nums := []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}
		got := containsDuplicate(nums)
		want := true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
