package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func greetings(ctx context.Context, gopher string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for {
			// limit the rate of the output
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

			select {
			case <-ctx.Done():
				return
			case c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher):
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	ctx, cancelFn := context.WithCancel(context.Background())

	c1 := greetings(ctx, "Goffredo")
	c2 := greetings(ctx, "Golia")
	c3 := greetings(ctx, "Gaetano")

	stop := time.After(5 * time.Second)
	for {
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg)
		case msg := <-c3:
			fmt.Println(msg)
		case <-stop:
			cancelFn()
			return
		}
	}
}
