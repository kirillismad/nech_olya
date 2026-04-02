package main

func UniqueNumber(arr []int) bool {
	freq := make(map[int]int)
	for _, num := range arr {
		freq[num]++
	}

	seen := make(map[int]bool)
	for _, count := range freq {
		if seen[count] {
			return false
		}
		seen[count] = true
	}
	return true
}
