package every

import (
	"fmt"
)

type Animal interface {
	Sound() string
}

type Dog struct{ Name string }

func (d *Dog) Sound() string {
	return "Woof"
}
func (d *Dog) Fetch() string {
	return d.Name + " bring the ball"
}

type Cat struct{ Name string }

func (c *Cat) Sound() string {
	return "Meow"
}

func (c *Cat) Purr() string {
	return c.Name + " purr"
}

func Interact(a Animal) {
	if val, ok := a.(*Dog); ok {
		fmt.Println(val.Fetch())
		return
	}
	if val, ok := a.(*Cat); ok {
		fmt.Println(val.Purr())
		return
	} else {
		panic("unknown type")
	}
}

func Interact2(a Animal) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	Interact(a)
}
