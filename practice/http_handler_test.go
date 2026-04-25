package practice

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	t.Log("server url:", server.URL)
	select {}

}
