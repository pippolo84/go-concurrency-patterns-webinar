package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Message holds data produced by our service
type Message struct {
	payload string
}

// Service is the descriptor of the service
type Service struct {
	events chan Message
	errors chan error

	// additional status
}

// NewService creates a new Service
func NewService() (*Service, error) {
	Service := &Service{
		events: make(chan Message),
		errors: make(chan error),
	}

	return Service, nil
}

// Run starts the service and the production of messages
func (a *Service) Run(ctx context.Context, wg *sync.WaitGroup) (<-chan Message, <-chan error) {
	wg.Add(1)

	go func() {
		defer func() {
			// clean up the service during shutdown
			close(a.events)
			close(a.errors)

			wg.Done()
		}()

		for {
			var (
				msg Message
				err error
			)

			// simulate a job that takes some time and may produce errors
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			if rand.Intn(2) != 0 {
				msg, err = Message{"payload"}, nil
			} else {
				msg, err = Message{}, errors.New("error")
			}

			// send back results
			select {
			case <-ctx.Done():
				return
			case a.events <- msg:
			case a.errors <- err:
			}
		}
	}()

	return a.events, a.errors
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Create the service
	Service, err := NewService()
	if err != nil {
		panic(err)
	}

	// Start the service
	events, errors := Service.Run(ctx, &wg)

	// Consume events
	for {
		select {
		case <-stop:
			// signal cancellation
			cancel()
			// wait for all goroutine to complete shutdown
			wg.Wait()

			return
		case ev := <-events:
			fmt.Printf("received event: %v\n", ev)
		case err := <-errors:
			fmt.Printf("received error: %v\n", err)
		}
	}
}
