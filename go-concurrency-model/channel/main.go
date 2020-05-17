package main

import "fmt"

func greetings(gopher string) <-chan string {
	c := make(chan string)

	go func() {
		c <- fmt.Sprintf("Hello, I'm %s, nice to meet you!", gopher)
	}()

	return c
}

func main() {
	fmt.Println(<-greetings("Goffredo"))
	fmt.Println(<-greetings("Golia"))
	fmt.Println(<-greetings("Gottardo"))
}
