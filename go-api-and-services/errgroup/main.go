package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

var urls = []string{
	"http://httpbin.org/range/1024?duration=4&chunk_size=256",
	"http://httpbin.org/range/2048?duration=8&chunk_size=256",
	"http://httpbin.org/range/512?duration=2&chunk_size=256",
}

func slowHandler(w http.ResponseWriter, req *http.Request) {
	// create a timeout context
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// create an errgroup with context derived from the previous one
	g, ctx := errgroup.WithContext(timeoutCtx)

	// requests 256 bytes each second (total 1024 bytes)
	for _, url := range urls {
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		url := url

		// launch a goroutine to fetch urls in parallel
		// each error will be correctly propagated to the caller
		// the timeout context will avoid blocking for too long
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return ctx.Err()
			}

			// send request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return ctx.Err()
			}
			defer resp.Body.Close()

			// read body and discard it
			if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
				return ctx.Err()
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Printf("error occured: %v\n", err)
	}

	w.Write([]byte("success"))
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
