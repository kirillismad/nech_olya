package leetcode

// O(n^2)-по времени
// O(1)-по памяти
func isAnagram(s, t string) bool {
	var count int

	if len(s) > len(t) {
		count = len(s)
	} else {
		count = len(t)
	}

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(t); j++ {
			if s[i] == t[j] {
				count--
			}
		}
	}

	if count > 0 {
		return false
	}
	return true
}
