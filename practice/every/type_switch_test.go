package every

import (
	"fmt"
	"testing"
)

func TestDescribe(t *testing.T) {
	var a interface{}
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	str := []string{"hello", "world"}
	fmt.Println(Describe(55))
	fmt.Println(Describe("hello"))
	fmt.Println(Describe(true))
	fmt.Println(Describe(12.0))
	fmt.Println(Describe(a))
	fmt.Println(Describe(nums))
	fmt.Println(Describe(str))
}
