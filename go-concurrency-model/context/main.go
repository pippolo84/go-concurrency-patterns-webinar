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
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher)
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			}
		}
	}()
	return c
}

func f() {
	rand.Seed(time.Now().Unix())
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)

	c1 := greetings(ctx, "Goffredo")
	c2 := greetings(ctx, "Golia")
	stop := time.After(time.Duration(5) * time.Second)
	for {
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg)
		case <-stop:
			cancelFn()
			return
		}
	}
}

func main() {
	f()
}
