package practice

import (
	"errors"
	"log"
)

var (
	ErrName  = errors.New("empty name")
	ErrQty   = errors.New("qty cannot be 0 or less")
	ErrPrice = errors.New("price annot be 0")
)

func ValidateOrder(name string, qty int, price float64) error {
	var errs []error
	if name == "" {
		errs = append(errs, errors.New("empty name"))
	}
	if qty <= 0 {
		errs = append(errs, errors.New("qty cannot be 0 or less"))
	}
	if price < 0 {
		errs = append(errs, errors.New("price annot be 0"))
	}
	return errors.Join(errs...)
}
func errorsJoin(name string, qty int, price float64) {
	err := ValidateOrder(name, qty, price)
	if err != nil {
		if errors.Is(err, ErrName) {
			log.Printf("ErrName== %s", err)
		}
		if errors.Is(err, ErrQty) {
			log.Printf("ErrQty==%s", err)
		}
		if errors.Is(err, ErrPrice) {
			log.Printf("ErrPrice==%s", err)
		}
	}
}
