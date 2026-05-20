package practice

import "testing"

func TestGetRequest(t *testing.T){
	got:=getRequest()
	want:="200 OK"
	if got!=want{
		t.Error("expected code 200")
		return
	}
}