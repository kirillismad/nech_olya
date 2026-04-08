package leetcode

// O(n)-по времени
// O(1)-по памяти
func twoSum(nums []int, target int) []int {
	elemNums := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		elem := nums[i]
		twoElem := target - elem
		idElem, ok := elemNums[twoElem]
		if ok {
			return []int{idElem, i}
		}
		elemNums[elem] = i
	}
	return nil
}
