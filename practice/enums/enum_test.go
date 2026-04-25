package enums

import "testing"

func TestEnum(t *testing.T) {
	var got Number = NumberNine
	want := 9
	if int(got) != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
