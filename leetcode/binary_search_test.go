package main

import "testing"

func TestBinarySearch(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	target := 1
	got := BinarySearch(nums, target)
	want := 0
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
