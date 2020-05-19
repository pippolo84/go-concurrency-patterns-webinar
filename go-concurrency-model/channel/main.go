package main

import "fmt"

func greetings(gophers ...string) <-chan string {
	c := make(chan string)

	go func() {
		defer close(c)
		for _, gopher := range gophers {
			c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher)
		}
	}()

	return c
}

func main() {
	gophers := []string{"Gaetano", "Goffredo", "Golia", "Gottardo"}
	for message := range greetings(gophers...) {
		fmt.Println(message)
	}
}
