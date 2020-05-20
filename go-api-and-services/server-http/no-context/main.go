package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, req *http.Request) {
	// requests 256 bytes each second (total 1024 bytes)
	req, err := http.NewRequest("GET", "http://httpbin.org/range/1024?duration=4&chunk_size=256", nil)
	if err != nil {
		panic(err)
	}

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// read body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
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
