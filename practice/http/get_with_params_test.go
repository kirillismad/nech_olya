package http

import "testing"

func TestGetWithParams(t *testing.T) {
	code,err:=getWithParams()
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