package main

func BinarySerch(nums []int, target int) int {
	mid := len(nums)/2 - 1
	ind := -1
	if target <= nums[mid] {
		for i := 0; i < len(nums); i++ {
			if target == nums[i] {
				ind = i
			}
		}
	} else if target >= nums[mid] {
		for i := len(nums) / 2; i < len(nums); i++ {
			if target == nums[i] {
				ind = i
			}
		}
	}
	return ind
}
