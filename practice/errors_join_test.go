package practice

import (
	"errors"
	"testing"
)

func TestErrorJoin(t *testing.T) {
	t.Run("tc 1", func(t *testing.T) {
		var name string
		var qty int
		price := -1.1
		err := ValidateOrder(name, qty, price)
		if err != nil {
			if errors.Is(err, ErrName) && errors.Is(err, ErrQty) && errors.Is(err, ErrPrice) {
				t.Logf("err== %v, ", err)
			} else {
				t.Errorf("unknown error:%v ", err)
				return
			}
		}
		if err == nil {
			t.Error("ошибок нет")
			return
		}
	})
	t.Run("tc 2", func(t *testing.T) {
		var name string
		qty := 5
		price := 100.2
		err := ValidateOrder(name, qty, price)
		if err != nil {
			if errors.Is(err, ErrName) {
				t.Logf("ErrName== %s", err)
			} else {
				t.Errorf("unknown error: %v", err)
				return
			}
		}
		if err == nil {
			t.Error("ошибок нет")
			return
		}
	})

	t.Run("tc 3", func(t *testing.T) {
		name := "milk"
		qty := 5
		price := 100.2
		err := ValidateOrder(name, qty, price)
		if err != nil {
			t.Errorf("this result is not expected, error, %v", err)
			return
		}
		if err == nil {
			t.Log("no errors")
			return
		}
	})
}
