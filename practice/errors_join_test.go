package practice

import "testing"

func TestErrorJoin(t *testing.T) {
	t.Run("tc 1", func(t *testing.T) {
		var name string
		var qty int
		price := -1.1
		errorsJoin(name, qty, price)
	})
	t.Run("tc 2", func(t *testing.T) {
		var name string
		qty := 5
		price := 100.2
		errorsJoin(name, qty, price)
	})

	t.Run("tc 3", func(t *testing.T) {
		name := "milk"
		qty := 5
		price := 100.2
		errorsJoin(name, qty, price)
	})
}
