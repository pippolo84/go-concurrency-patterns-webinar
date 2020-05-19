package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Task represents a task that will take duration to execute
type Task struct {
	ID       int
	duration time.Duration
}

// NewTask creates a new task to be executed
func NewTask(ID int) Task {
	return Task{
		ID:       ID,
		duration: time.Duration(time.Duration(rand.Intn(250)) * time.Millisecond),
	}
}

// Complete executes the task
func (t *Task) Complete() {
	time.Sleep(t.duration)
}

// String returns a human-readable representation of the task
func (t Task) String() string {
	return fmt.Sprintf("task %d", t.ID)
}

// Generator turns a slice of tasks in a stream of tasks
func Generator(tasks ...Task) <-chan Task {
	stream := make(chan Task)

	go func() {
		defer close(stream)

		for _, t := range tasks {
			stream <- t
		}
	}()

	return stream
}

// Stage executes a step of a pipeline on each task from input,
// giving back the result in the returned stream
func Stage(step string, input <-chan Task) <-chan Task {
	output := make(chan Task)

	go func() {
		defer close(output)

		for t := range input {
			// ...work on task t...
			fmt.Printf("applying step %q to %v\n", step, t)
			t.Complete()

			output <- t
		}
	}()

	return output
}

func main() {
	ntask := 10

	// our tasks to complete
	tasks := make([]Task, ntask)
	for i := 0; i < ntask; i++ {
		tasks[i] = NewTask(i)
	}

	// generate a stream to feed the pipeline
	stream := Generator(tasks...)
	// create the pipeline
	pipeline := Stage("third", Stage("second", Stage("first", stream)))
	// drain the pipeline
	for t := range pipeline {
		fmt.Printf("task completed: %v\n", t)
	}
}
