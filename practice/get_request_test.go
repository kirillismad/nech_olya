package practice

import "testing"

func TestGetRequest(t *testing.T) {
	code, err := getRequest()
	if code != 200 {
		t.Error("expected code 200")
		return
	}
	if err != nil {
		t.Error("error was not expected")
		return
	}
	t.Log(code)
}
