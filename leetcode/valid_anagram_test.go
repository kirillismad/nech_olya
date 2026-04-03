package leetcode

import "testing"

func TestValidAnagram(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		got := isAnagram("anagram", "nagaram")
		want := true

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		got := isAnagram("rat", "car")
		want := false

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
