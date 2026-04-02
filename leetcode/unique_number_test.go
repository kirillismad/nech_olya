package main

import "testing"

func TestUniqueNumber(t *testing.T) {
	arr := []int{1, 2, 2, 1, 1, 0}
	got := UniqueNumber(arr)
	want := true
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
