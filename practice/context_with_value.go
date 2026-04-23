package practice

import (
	"context"
	"fmt"
)

type key int

const (
	requestIDKey key = 0
	falseIDKey       = 1
)

func contextWithValue() {
	ctx := context.WithValue(context.Background(), requestIDKey, 123)
	requestkey := getItValue(ctx, requestIDKey)
	fmt.Printf("true key: %v\n", requestkey)
	falseKey := getItValue(ctx, falseIDKey)
	fmt.Printf("false key: %v\n", falseKey)

}

func getItValue(ctx context.Context, k key) (val int) {
	requestID := ctx.Value(k)
	if requestID == nil {
		return 0
	}
	return (requestID).(int)
}
