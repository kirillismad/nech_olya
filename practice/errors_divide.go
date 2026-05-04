package practice

import (
	"errors"
	"fmt"
	"strconv"
)

func Divide(a, b float64) (float64, error) {
	if b == 0.0 {
		return 0.0, errors.New("cannot be divided by 0")
	}
	return a / b, nil
}

func ParseAge(s string) (int, error) {
	age, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to number, error: %w", err)
	}
	return age, nil
}
