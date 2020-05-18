package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

type message struct{}

func producer(ctx context.Context, wg *sync.WaitGroup) <-chan message {
	c := make(chan message)

	go func() {
		defer func() {
			close(c)
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// produce messages
				c <- message{}

				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	return c
}

func consumer(ctx context.Context, wg *sync.WaitGroup, c <-chan message) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-c:
				// consume messages
				fmt.Println(msg)
			}
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	wg.Add(2)
	c := producer(ctx, &wg)
	consumer(ctx, &wg, c)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	cancelFn()

	wg.Wait()
}
