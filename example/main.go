package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CyCoreSystems/errchan"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ec := errchan.New()
	defer ec.Close()

	ec.Wrap(func() error {
		return startA(ctx)
	})

	ec.Wrap(func() error {
		return startB(ctx)
	})

	select {
	case <-ctx.Done():
	case err := <-ec.Next():
		fmt.Println("Got error:", err)
	}
}

func startA(ctx context.Context) error {
	time.Sleep(10 * time.Second)
	return errors.New("Timeout after 10 seconds")
}

func startB(ctx context.Context) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Pong")
	})
	return http.ListenAndServe(":8080", nil)
}
