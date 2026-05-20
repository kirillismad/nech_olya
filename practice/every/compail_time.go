package every

import "fmt"

type Describer interface {
	Describe() string
}

type User struct {
	Name string
}

func (u *User) Describe() string {
	return u.Name
}

type Order struct {
	Address string
}

func (o *Order) Describe() string {
	return o.Address
}

type Product struct {
	Name string
}

func (p *Product) Describe() string {
	return p.Name
}

var _ Describer = (*User)(nil)
var _ Describer = (*Order)(nil)
var _ Describer = (*Product)(nil)

func PrintAll(items []Describer) {
	for _, item := range items {
		fmt.Println(item.Describe())
	}
}

func compailTime() {
	items := []Describer{
		&User{Name: "Olya"},
		&Order{Address: "Lenin street"},
		&Product{Name: "Bread"}}
	PrintAll(items)
}
