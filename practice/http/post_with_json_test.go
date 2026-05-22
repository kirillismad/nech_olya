package http

import "testing"

func TestPostWithJSON(t *testing.T) {
	code,err:=postWithJson()
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
