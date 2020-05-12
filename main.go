package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

type Message struct{}

type Actor struct {
	Events chan<- Message
	Errors chan<- error

	// additional status
}

func NewActor() (*Actor, error) {
	actor := &Actor{
		Events: make(chan Message),
		Errors: make(chan error),
	}

	return actor, nil
}

func (a *Actor) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// do work and send results to Events and Errors channels
			}
		}
	}()
}

func main() {
	actor, err := NewActor()
	if err != nil {
		panic(err)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	actor.Run(ctx, &wg)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	cancelFn()

	// wait for all goroutine to complete shutdown
	wg.Wait()
}
