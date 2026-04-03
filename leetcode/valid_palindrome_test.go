package leetcode

import "testing"

func TestValidPalindrome(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		got := isPalindrome("A man, a plan, a canal: Panama")
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		got := isPalindrome("race a car")
		want := false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("tc 3:", func(t *testing.T) {
		got := isPalindrome(" ")
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
