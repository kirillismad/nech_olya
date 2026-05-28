package every

import (
	"fmt"
	"testing"
)

func TestRecursion(t *testing.T) {
	testsSum := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{5, 15},
		{7, 28},
	}
	for _, test := range testsSum {
		result := SumToN(test.n)
		if result != test.expected {
			t.Fatalf(
				"expected %d, got %d",
				test.expected,
				result,
			)
		}
	}

	testsFactorial := []struct {
		n        int
		expected int
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{7, 5040},
	}
	for _, test := range testsFactorial {
		t.Run(fmt.Sprintf("Factorial(%d)", test.n), func(t *testing.T) {
			result := Factorial(test.n)
			if result != test.expected {
				t.Fatalf(
					"expected %d, got %d",
					test.expected,
					result,
				)
			}
		})
	}
}
