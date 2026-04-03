package leetcode

import "testing"

func TestTwoSum(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		nums := []int{2, 7, 11, 15}
		target := 9
		got := twoSum(nums, target)
		want := []int{0, 1}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		nums := []int{3, 2, 4}
		target := 6
		got := twoSum(nums, target)
		want := []int{1, 2}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})
	t.Run("tc 3:", func(t *testing.T) {
		nums := []int{3, 3}
		target := 6
		got := twoSum(nums, target)
		want := []int{0, 1}
		for i := 0; i < len(got); i++ {
			if got[i] != want[i] {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})

}
