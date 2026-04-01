package main

import "testing"

func TestBinarySerch(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	target := 1
	got := BinarySerch(nums, target)
	want := 0
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
