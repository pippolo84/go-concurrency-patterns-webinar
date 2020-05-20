package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Item is a generic unit of work
type Item struct{}

// Callback is the type of operation to be executed on a Item
type Callback func(Item)

// Service fetches elements asynchronously
type Service struct {
	sync.Mutex
	cb Callback
}

// NewService returns a newly allocated Service
func NewService() *Service {
	return &Service{}
}

// SetCallback sets a user-defined callback to be called on each fetched item
func (s *Service) SetCallback(callback func(Item)) {
	s.Lock()
	defer s.Unlock()

	s.cb = callback
}

// Run fetches elements asynchronously and execute the callback on each of them
func (s *Service) Run() {
	go func() {
		for {
			// simulate a long operation to fetch an item
			time.Sleep(time.Duration(500) * time.Millisecond)

			// load callback
			s.Lock()
			cb := s.cb
			s.Unlock()

			// execute callback on fetched item
			if cb != nil {
				cb(Item{})
			}

		}
	}()
}

// Operation is an example of a callback to be applied to an item
func Operation(i Item) {
	fmt.Printf("executing callback on item: %v\n", i)
}

func main() {
	service := NewService()
	service.SetCallback(Operation)
	service.Run()

	for {
	}
}
