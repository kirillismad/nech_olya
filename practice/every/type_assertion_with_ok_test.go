package every

import "testing"

func TestTypeAssertion(t *testing.T) {
	dog := Dog{Name: "Bob"}
	cat := Cat{Name: "Vasya"}
	Interact2(&dog)
	Interact2(&cat)
}
