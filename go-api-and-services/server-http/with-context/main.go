package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// requests 256 bytes each second (total 1024 bytes)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://httpbin.org/range/1024?duration=4&chunk_size=256", nil)
	if err != nil {
		panic(err)
	}

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read body
	if _, err := io.Copy(w, resp.Body); err != nil {
		// check for context timeout
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("timeout!")
			return
		}

		panic(err)
	}
}

func main() {
	srv := http.Server{
		Addr:         ":8080",
		WriteTimeout: 1 * time.Second,
		Handler:      http.HandlerFunc(slowHandler),
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
