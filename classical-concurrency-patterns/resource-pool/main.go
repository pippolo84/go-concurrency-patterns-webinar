package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// a counter to keep track of the created resources
var counter int32

// Resource is a mock for a generic resource
type Resource struct {
	id int32
}

// Pool manages a bounded number of equivalent resources
type Pool struct {
	free chan *Resource
	busy chan struct{}
}

// NewPool creates a Pool of bound resources
func NewPool(bound int) *Pool {
	return &Pool{
		free: make(chan *Resource, bound),
		busy: make(chan struct{}, bound),
	}
}

// Get acquires a Resource
func (p *Pool) Get(ctx context.Context) *Resource {
	p.busy <- struct{}{}
	select {
	case <-ctx.Done():
		return nil
	case r := <-p.free:
		return r
	default:
		// lazily create a new resource
		resourceID := atomic.AddInt32(&counter, 1)
		return &Resource{
			id: resourceID,
		}
	}
}

// Put releases a Resource
func (p *Pool) Put(ctx context.Context, r *Resource) {
	select {
	case <-ctx.Done():
		return
	default:
		p.free <- r
		<-p.busy
	}
}

func worker(ctx context.Context, wg *sync.WaitGroup, p *Pool) {
	go func() {
		var r *Resource
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				if r != nil {
					p.Put(ctx, r)
					r = nil
				} else {
					r = p.Get(ctx)
					if r == nil {
						return
					}

					// use acquired resource
					fmt.Printf("using resource: %v\n", *r)
				}
			}
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	pool := NewPool(10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		worker(ctx, &wg, pool)
	}

	wg.Wait()
}
