package main

import "unicode"

func validPalindrome(s string) bool {
	str := cleanStrings(s)
	if str == "" {
		return true
	}

	for i := 0; i < len(str); i++ {
		for j := len(str) - 1; j >= 0; j-- {
			if str[i] == str[j] {
				return true
			}
		}
	}
	return false
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
