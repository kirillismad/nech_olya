package practice

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
		if err != nil {
			if val, ok := err.(*ValidationError); ok {
				if name == "" {
					t.Logf("Field: %s\n Message: %s", val.Field, val.Message)
				} else if age < 0 {
					t.Errorf("Field: %s\n Message: %s", val.Field, val.Message)
					return
				} else {
					t.Errorf("another error")
					return
				}

			}
		}
	})
	t.Run("tc 2:", func(t *testing.T) {
		name := "olya"
		age := 25
		err := ValidateUser(name, age)
		if err != nil {
			if val, ok := err.(*ValidationError); ok {
				if name == "" {
					t.Errorf("Field: %s\n Message: %s", val.Field, val.Message)
					return
				} else if age < 0 {
					t.Errorf("Field: %s\n Message: %s", val.Field, val.Message)
					return
				}

				t.Errorf("another error")
				return

			}
			if err == nil {
				t.Log("data entered correctly")
			}
		}
	})
	t.Run("tc 3:", func(t *testing.T) {
		name := "Alice"
		age := -10
		err := ValidateUser(name, age)
		if err != nil {
			if val, ok := err.(*ValidationError); ok {
				if name == "" {
					t.Errorf("Field: %s\n Message: %s", val.Field, val.Message)
				} else if age < 0 {
					t.Logf("Field: %s\n Message: %s", val.Field, val.Message)
					return
				} else {
					t.Errorf("another error")
					return
				}
			}
			if err == nil {
				t.Error("data entered correctly")
				return
			}
		}
	})

	t.Run("tc 4:", func(t *testing.T) {
		id := 0
		err := FindUserByID(id)
		if err != nil {
			if val, ok := err.(*NotFoundError); ok {
				t.Errorf("ID: %d", val.ID)
				return
			} else {
				t.Errorf("another error")
				return
			}
		}
		if err == nil {
			t.Log("data entered correctly")
		}
	})
	t.Run("tc 5:", func(t *testing.T) {
		id := 100
		err := FindUserByID(id)
		if err != nil {
			if val, ok := err.(*NotFoundError); ok {
				t.Logf("Error ID: %d", val.ID)
			} else {
				t.Errorf("another error")
				return
			}
		}
		if err == nil {
			t.Error("data entered correctly")
			return
		}
	})
}
