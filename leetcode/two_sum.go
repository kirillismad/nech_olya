package leetcode

// O(n^2)-по времени
// O(1)-по памяти
func twoSum(nums []int, target int) []int {
	s := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				s = append(s, i, j)
			}
		}
	}
	return s
}
