package practice

import (
	"errors"
)

var (
	ErrName  = errors.New("empty name")
	ErrQty   = errors.New("qty cannot be 0 or less")
	ErrPrice = errors.New("price annot be 0")
)

func ValidateOrder(name string, qty int, price float64) error {
	var errs []error
	if name == "" {
		errs = append(errs, ErrName)
	}
	if qty <= 0 {
		errs = append(errs, ErrQty)
	}
	if price < 0 {
		errs = append(errs, ErrPrice)
	}
	return errors.Join(errs...)
}
