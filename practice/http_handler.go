package practice

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-time.After(time.Second):
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		fmt.Println(w, "request completed")
	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("Запрос отменен", err)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}
}

func httpHandler() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
