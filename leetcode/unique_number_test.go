package leetcode

import "testing"

func TestUniqueNumber(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		arr := []int{1, 2, 2, 1, 1, 3}
		got := uniqueOccurrences(arr)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		arr := []int{1, 2}
		got := uniqueOccurrences(arr)
		want := false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("tc 3:", func(t *testing.T) {
		arr := []int{-3, 0, 1, -3, 1, 1, 1, -3, 10, 0}
		got := uniqueOccurrences(arr)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
