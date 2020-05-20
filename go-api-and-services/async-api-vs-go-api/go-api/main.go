package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Item is a generic unit of work
type Item struct{}

// Service is the descriptor of the service
type Service struct{}

// NewService creates a new Service
func NewService() *Service {
	return &Service{}
}

// Run fetches elements asynchronously and send them back through the returned channel
func (s *Service) Run() <-chan Item {
	items := make(chan Item)
	go func() {
		defer close(items)

		for {
			// simulate a long operation to fetch an item
			time.Sleep(time.Duration(500) * time.Millisecond)
			items <- Item{}
		}
	}()

	return items
}

// Operation is an example of a callback to be applied to an item
func Operation(i Item) {
	fmt.Printf("executing operation on item: %v\n", i)
}

func main() {
	service := NewService()
	items := service.Run()
	go func() {
		for item := range items {
			Operation(item)
		}
	}()

	for {
	}
}
