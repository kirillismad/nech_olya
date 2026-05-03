package practice

import "testing"

func TestCustomType(t *testing.T) {
	t.Run("tc 1:", func(t *testing.T) {
		name := ""
		age := -1
		err := ValidateUser(name, age)
		if err == nil {
			t.Log("data entered correctly")
		}
		if err != nil {
			if val, ok := err.(*ValidationError); ok {
				t.Errorf("Field: %s", val.Field)
				t.Errorf("Message: %s", val.Message)
			} else {
				t.Errorf("another error")
			}
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		name := "olya"
		age := 25
		err := ValidateUser(name, age)
		if err == nil {
			t.Log("data entered correctly")
		}
		if err != nil {
			if val, ok := err.(*ValidationError); ok {
				t.Errorf("Field: %s", val.Field)
				t.Errorf("Message: %s", val.Message)
			} else {
				t.Errorf("another error")
			}
		}
	})

	t.Run("tc 3:", func(t *testing.T) {
		id := 123
		err := FindUserByID(id)
		if err == nil {
			t.Log("data entered correctly")
		}
		if err != nil {
			if val, ok := err.(*NotFoundError); ok {
				t.Errorf("ID: %d", val.ID)
			} else {
				t.Errorf("another error")
			}
		}
	})
	t.Run("tc 4:", func(t *testing.T) {
		id := 100
		err := FindUserByID(id)
		if err == nil {
			t.Log("data entered correctly")
		}
		if err != nil {
			if val, ok := err.(*NotFoundError); ok {
				t.Errorf("Error ID: %d", val.ID)
			} else {
				t.Errorf("another error")
			}
		}
	})
}
