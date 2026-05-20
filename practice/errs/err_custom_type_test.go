package errs

import "testing"

func TestCustomType(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		name := ""
		age := 25
		err := ValidateUser(name, age)
		if err == nil {
			t.Error("data entered correctly")
			return
		}
		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Errorf("unexpected error type: %T", err)
			return
		}
		if validationErr.Field != "name" {
			t.Errorf("unexpected field: %s", validationErr.Field)
			return
		}
		if validationErr.Message != "cannot be empty" {
			t.Errorf("unexpected message: %s", validationErr.Message)
			return
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		name := "olya"
		age := 25
		err := ValidateUser(name, age)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
	t.Run("tc 3: error: negative age", func(t *testing.T) {
		name := "Alice"
		age := -10
		err := ValidateUser(name, age)
		if err == nil {
			t.Error("data entered correctly")
			return
		}
		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Errorf("unexpected error type: %T", err)
			return
		}
		if validationErr.Field != "age" {
			t.Errorf("unexpected field: %s", validationErr.Field)
			return
		}
		if validationErr.Message != "cannot be negative" {
			t.Errorf("unexpected message: %s", validationErr.Message)
			return
		}
	})

	t.Run("tc 4: success", func(t *testing.T) {
		id := 0
		err := FindUserByID(id)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
	t.Run("tc 5: error: user not found", func(t *testing.T) {
		id := 100
		err := FindUserByID(id)
		if err == nil {
			t.Error("user found")
			return
		}
		notFoundErr, ok := err.(*NotFoundError)
		if !ok {
			t.Errorf("unexpected error type: %T", err)
			return
		}
		if notFoundErr.ID != id {
			t.Errorf("unexpected ID: %d", notFoundErr.ID)
			return
		}
	})
}
