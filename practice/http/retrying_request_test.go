package http

import "testing"

func TestRetry(t *testing.T){
	t.Run("tc 1:",func(t *testing.T) {
		retryRequest("https://httpbin.org/status/500")
	})
}