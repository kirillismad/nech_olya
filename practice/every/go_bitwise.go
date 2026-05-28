package every

import (
	"fmt"
)

func BitAnd(a, b uint8) uint8 {
	var c uint8 = a & b
	fmt.Printf("c: %d, %b\n", c, c)
	return a & b
}

func BitOr(a, b uint8) uint8 {
	var c uint8 = a | b
	fmt.Printf("c: %d, %b\n", c, c)
	return a | b
}

func BitXor(a, b uint8) uint8 {
	var c uint8 = a ^ b
	fmt.Printf("c: %d, %b\n", c, c)
	return a ^ b
}

func BitNot(a, b uint8) uint8 {
	var c uint8 = a &^ b
	fmt.Printf("c: %d, %b\n", c, c)
	return a &^ b
}
