package http

import (
	"testing"
)

func TestRetryRequst(t *testing.T) {
	_, err := retry()
	if err == nil {
		t.Error("error expected")
		return
	}
	status500 := "500"
	if err.Error() != status500 {
		t.Error("error expected 500!")
		return
	}
	t.Log("the test was successful")
}
