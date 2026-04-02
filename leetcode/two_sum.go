package main

func twoSum(nums []int, target int) (s []int) {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				s = append(s, i, j)
			}
		}
	}
	return s
}
