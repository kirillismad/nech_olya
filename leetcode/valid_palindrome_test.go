package main

import "testing"

func TestValidPalindrome(t *testing.T) {
	got := validPalindrome("A man, a plan, a canal: Panama")
	want := true
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
