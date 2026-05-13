package practice

import (
	"fmt"
	"testing"
)

func TestDivide(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		a := 2.0
		b := 0.0
		_, err := Divide(a, b)
		if err != nil {
			fmt.Println(err)
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		a := 3.0
		b := 1.5
		age, err := Divide(a, b)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(age)
	})

	t.Run("tc 3:", func(t *testing.T) {
		str:="25"
		age, err := ParseAge(str)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(age)
	})

	t.Run("tc 4:", func(t *testing.T) {
		str:="двадцать"
		age, err := ParseAge(str)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(age)
	})
}
