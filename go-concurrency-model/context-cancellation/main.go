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
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(rand.Intn(3)) * time.Second):
				c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher)
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	ctx, cancelFn := context.WithCancel(context.Background())

	c := greetings(ctx, "Goffredo")
	stop := time.After(5 * time.Second)
	for {
		select {
		case msg := <-c:
			fmt.Println(msg)
		case <-stop:
			cancelFn()
			return
		}
	}
}
