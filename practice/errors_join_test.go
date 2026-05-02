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
		if err == nil {
			t.Log("ошибок нет")
			return
		}
		if err != nil {
			if errors.Is(err, ErrName) {
				t.Logf("ErrName == %v", err)
			}
			if errors.Is(err, ErrQty) {
				t.Logf("ErrQty ==%v", err)
			}
			if errors.Is(err, ErrPrice) {
				t.Logf("ErrPrice ==%v", err)
			}
		}
	})
	t.Run("tc 2", func(t *testing.T) {
		var name string
		qty := 5
		price := 100.2
		err := ValidateOrder(name, qty, price)
		if err == nil {
			t.Log("ошибок нет")
			return
		}
		if err != nil {
			if errors.Is(err, ErrName) {
				t.Logf("ErrName== %s", err)
			}
			if errors.Is(err, ErrQty) {
				t.Logf("ErrQty==%s", err)
			}
			if errors.Is(err, ErrPrice) {
				t.Logf("ErrPrice==%s", err)
			}
		}
	})

	t.Run("tc 3", func(t *testing.T) {
		name := "milk"
		qty := 5
		price := 100.2
		err := ValidateOrder(name, qty, price)
		if err == nil {
			t.Log("ошибок нет")
			return
		}
		if err != nil {
			if errors.Is(err, ErrName) {
				t.Logf("ErrName== %s", err)
			}
			if errors.Is(err, ErrQty) {
				t.Logf("ErrQty==%s", err)
			}
			if errors.Is(err, ErrPrice) {
				t.Logf("ErrPrice==%s", err)
			}
		}
	})
}
