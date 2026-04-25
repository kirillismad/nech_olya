package context

import (
	"context"
	"fmt"
)

type key int

const (
	requestIDKey key = 0
	falseIDKey   key = 1
)

func ContextWithValue() {
	ctx := context.WithValue(context.Background(), requestIDKey, 123)
	requestkey := GetItValue(ctx, requestIDKey)
	fmt.Printf("true key: %v\n", requestkey)
	falseKey := GetItValue(ctx, falseIDKey)
	fmt.Printf("false key: %v\n", falseKey)

}

func GetItValue(ctx context.Context, k key) (val int) {
	requestID := ctx.Value(k)
	if requestID == nil {
		return 0
	}
	return (requestID).(int)
}
