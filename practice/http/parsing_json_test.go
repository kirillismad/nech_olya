package http

import "testing"

func TestParsingJSON(t *testing.T) {
	t.Run("tc 1: success", func(t *testing.T) {
		code, err := parsingJSONUnmarshal()
		if code != 200 {
			t.Error("expected code 200")
			return
		}
		if err != nil {
			t.Error("error was not expected")
			return
		}
		t.Log(code)
	})
	t.Run("tc 2: success", func(t *testing.T) {
		code, err := parsingJSONDecoder()
		if code != 200 {
			t.Error("expected code 200")
			return
		}
		if err != nil {
			t.Error("error was not expected")
			return
		}
		t.Log(code)
	})
}
