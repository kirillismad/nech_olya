package main

import "testing"

func TestValidAnagram(t *testing.T) {
	got := ValidAnagram("listen", "silent")
	want := true

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
