package leetcode

// O(n)-по времени
// O(n)-по памяти

func isAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	count := make(map[rune]int)

	for _, v := range s {
		count[v]++
	}

	for _, v := range t {
		count[v]--
		if count[v] < 0 {
			return false
		}
	}
	return true

}
