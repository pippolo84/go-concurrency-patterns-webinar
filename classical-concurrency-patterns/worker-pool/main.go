package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Task represents a task that will take duration to execute
type Task struct {
	ID       int
	duration time.Duration
}

// NewTask creates a new task to be executed
func NewTask(ID int) Task {
	return Task{
		ID:       ID,
		duration: time.Duration(rand.Intn(500)) * time.Millisecond,
	}
}

// Complete executes the task
func (t *Task) Complete() {
	time.Sleep(t.duration)
	fmt.Printf("task %d completed\n", t.ID)
}

// Token represents a place card that must be held to do something
type Token struct{}

// WorkerPool is a pool with an upper bound to the number of workers
type WorkerPool struct {
	busy chan Token
}

// NewWorkerPool create a worker pool with up to bound workers
func NewWorkerPool(bound int) (*WorkerPool, error) {
	return &WorkerPool{
		busy: make(chan Token, bound),
	}, nil
}

// Do launch a goroutine to complete tasks concurrently
func (wp *WorkerPool) Do(tasks []Task) {
	for _, task := range tasks {
		// acquire a token
		wp.busy <- Token{}

		go func(task Task) {
			// do work concurrently
			task.Complete()

			// release the token
			<-wp.busy
		}(task)
	}
}

func main() {
	tasks := make([]Task, 100)
	for i := 0; i < 100; i++ {
		tasks[i] = NewTask(i)
	}

	wp, err := NewWorkerPool(10)
	if err != nil {
		panic(err)
	}

	wp.Do(tasks)
}
