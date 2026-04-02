package main

import "testing"

func TestContainsDuplicate(t *testing.T) {
	nums := []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}
	got := containsDuplicate(nums)
	want := true
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
