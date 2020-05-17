package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func greetings(ctx context.Context, wg *sync.WaitGroup, gopher string) <-chan string {
	c := make(chan string)

	go func() {
		defer func() {
			close(c)
			wg.Done()
		}()

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
	var wg sync.WaitGroup

	ctx, cancelFn := context.WithCancel(context.Background())

	wg.Add(1)
	c := greetings(ctx, &wg, "Goffredo")

	stop := time.After(5 * time.Second)
	for {
		select {
		case msg := <-c:
			fmt.Println(msg)
		case <-stop:
			cancelFn()
			wg.Wait()
			return
		}
	}
}
