package practice

import (
	"fmt"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct{ Radius float64 }

func (c *Circle) Area() float64 {
	s := 3.14 * c.Radius * c.Radius
	return s
}

func (c *Circle) Perimeter() float64 {
	p := 2.0 * 3.14 * c.Radius
	return p
}

type Rectangle struct{ Width, Height float64 }

func (r *Rectangle) Area() float64 {
	s := r.Height * r.Width
	return s
}

func (r *Rectangle) Perimeter() float64 {
	p := 2 * (r.Height + r.Width)
	return p
}

func PrintInfo(s Shape) {
	fmt.Printf("Площадь фигуры-%v переиметр фигуры-%v\n", s.Area(), s.Perimeter())
}

func interfaseWithTwoType() {
	shape := []Shape{
		&Circle{Radius: 12.2},
		&Rectangle{Width: 5.0, Height: 3.0},
	}
	for i, _ := range shape {
		PrintInfo(shape[i])
	}
}
