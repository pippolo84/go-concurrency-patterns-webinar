package main

import (
	"fmt"
	"time"
)

func greetings() {
	fmt.Println("Hello, Go world!")
}

func main() {
	go greetings()

	// just wait a little bit
	time.Sleep(1 * time.Second)
}
