package every

func SumToN(n int) int {
	if n == 0 {
		return 0
	}
	return SumToN(n-1) + n
}

func Factorial(n int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return 1
	}
	return Factorial(n-1) * n
}
