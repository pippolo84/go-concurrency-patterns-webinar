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

	wg.Add(1)
	go func() {
		defer func() {
			fmt.Printf("%s is done!\n", gopher)

			close(c)
			wg.Done()
		}()

		for {
			// limit the rate of the output
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

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
	var wg sync.WaitGroup

	ctx, cancelFn := context.WithCancel(context.Background())

	c1 := greetings(ctx, &wg, "Goffredo")
	c2 := greetings(ctx, &wg, "Golia")
	c3 := greetings(ctx, &wg, "Gaetano")

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
			wg.Wait()
			return
		}
	}
}
