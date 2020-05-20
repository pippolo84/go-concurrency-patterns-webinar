package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func greetings(ctx context.Context, gopher string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for {
			// simulate a workload that takes a certain time to execute
			// but remain responsive listening for cancellation signal
			select {
			case <-ctx.Done():
				return
			case <-time.After(
				time.Duration(rand.Intn(500)) * time.Millisecond,
			):
			}

			// send the result back
			c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher)
		}
	}()
	return c
}

func main() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	c1 := greetings(ctx, "Goffredo")
	c2 := greetings(ctx, "Golia")
	c3 := greetings(ctx, "Gaetano")

	for {
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg)
		case msg := <-c3:
			fmt.Println(msg)
		case <-ctx.Done():
			return
		}
	}
}
