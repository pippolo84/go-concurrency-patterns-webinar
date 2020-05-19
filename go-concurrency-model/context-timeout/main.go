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
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	for msg := range greetings(ctx, "Goffredo") {
		fmt.Println(msg)
	}
}
