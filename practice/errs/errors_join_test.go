package errs

import (
	"errors"
	"testing"
)

func TestErrorJoin(t *testing.T) {
	t.Run("tc 1: error no name, qty zero, negative price", func(t *testing.T) {
		var name string
		var qty int
		price := -1.1
		err := ValidateOrder(name, qty, price)
		if err == nil {
			t.Error("error expected")
			return
		}
		if err != nil {
			if !errors.Is(err, ErrName) && !errors.Is(err, ErrQty) && !errors.Is(err, ErrPrice) {
				t.Errorf("unknown error:%v ", err)
				return
			}
			t.Logf("err== %v, ", err)
		}
	})
	t.Run("tc 2: error no name", func(t *testing.T) {
		var name string
		qty := 5
		price := 100.2
		err := ValidateOrder(name, qty, price)
		if err == nil {
			t.Error("ошибок нет")
			return
		}
		if err != nil {
			if !errors.Is(err, ErrName) {
				t.Errorf("unknown error: %v", err)
				return
			}
			t.Logf("ErrName== %s", err)
		}
	})

	t.Run("tc 3: success", func(t *testing.T) {
		name := "milk"
		qty := 5
		price := 100.2
		err := ValidateOrder(name, qty, price)
		if err != nil {
			t.Errorf("this result is not expected, error, %v", err)
			return
		}
		t.Log("no errors")
	})
}
