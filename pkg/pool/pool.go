package pool

import (
	"context"
	"fmt"
	"sync"
)

// This package contains implementation of a Thread Pool.

// ByteArray is an array of arrays (Slice of slices).
type ByteArray [][]byte

// Function -
type Function func([]byte)

// Coordinator is implementation of Thread Pool that uses one queue for deploying and executing tasks
type Coordinator struct {
	TaskQueue []Function
	DataQueue ByteArray
	CTX       context.Context
	mux       sync.Mutex
}

// CoordinatorInstance Global variable represents a single coordinator
var CoordinatorInstance = InitCoordinator()

// InitCoordinator initializes the coordinator
func InitCoordinator() *Coordinator {
	return &Coordinator{
		TaskQueue: make([]Function, 0),
		DataQueue: make([][]byte, 0),
	}
}

// Enqueue places a new task into the TaskQueue and returns its (TaskQueue's) length
func (c *Coordinator) Enqueue(fun func([]byte), data []byte) int {
	c.mux.Lock()
	c.TaskQueue = append(c.TaskQueue, fun)
	c.DataQueue = append(c.DataQueue, data)
	c.mux.Unlock()
	return len(c.TaskQueue)
}

// Dequeue removes one task and returns it to the caller
func (c *Coordinator) Dequeue() (func([]byte), []byte) {
	if len(c.TaskQueue) > 0 {
		c.mux.Lock()
		fun := c.TaskQueue[0]
		data := c.DataQueue[0]
		c.TaskQueue = c.TaskQueue[1:]
		c.DataQueue = c.DataQueue[1:]
		c.mux.Unlock()
		return fun, data
	}

	return nil, nil
}

// IsEmpty checks if coordinator queue is empty
func (c *Coordinator) IsEmpty() bool {
	return len(c.TaskQueue) == 0 || len(c.DataQueue) == 0
}

// Run runs in separate go thread
func (c *Coordinator) Run() {
	for {
		select {
		case <-c.CTX.Done():
			fmt.Println(" I am exiting ")
			break
		default:
			if !c.IsEmpty() {
				fun, data := c.Dequeue()
				fun(data)
			}
		}
	}
}
