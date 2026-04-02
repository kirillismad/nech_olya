package main

func containsDuplicate(nums []int) bool {
	repeat := make(map[int]int)
	for _, v := range nums {
		repeat[v]++
	}
	for _, v := range repeat {
		if v >= 2 {
			return true
		}
	}
	return false
}
