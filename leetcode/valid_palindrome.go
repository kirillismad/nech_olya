package leetcode

import "unicode"

// O(n)-по времени
// O(n)-по памяти
func isPalindrome(s string) bool {
	s = cleanStrings(s)
	if s == " " {
		return true
	}
	left, right := 0, len(s)-1

	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func cleanStrings(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, unicode.ToLower(r))
		}
	}
	return string(result)
}
