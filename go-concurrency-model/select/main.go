package main

import (
	"fmt"
	"math/rand"
	"time"
)

func greetings(gopher string) <-chan string {
	c := make(chan string)
	defer close(c)
	go func() {
		for {
			c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher)
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())

	c1 := greetings("Goffredo")
	c2 := greetings("Golia")
	for i := 0; i < 10; i++ {
		select {
		case msg := <-c1:
			fmt.Println(msg)
		case msg := <-c2:
			fmt.Println(msg)
		}
	}
}
