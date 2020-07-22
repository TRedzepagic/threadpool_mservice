package pool

import (
	"context"
	"sync"
)

// This package contains implementation of a Thread Pool.

// ByteArray is an array of arrays (Slice of slices).
type ByteArray [][]byte

// Function is a function that takes a s
type Function func([]byte)

// Wg represents our WaitGroup variable, used for graceful exit. Waits for current task in queue, then exits
// Ignores non-started tasks in queue.
var Wg sync.WaitGroup

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
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.TaskQueue) > 0 {
		fun := c.TaskQueue[0]
		data := c.DataQueue[0]
		c.TaskQueue = c.TaskQueue[1:]
		c.DataQueue = c.DataQueue[1:]
		return fun, data
	}

	return nil, nil
}

// IsEmpty checks if coordinator queue is empty
func (c *Coordinator) IsEmpty() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.TaskQueue) == 0 || len(c.DataQueue) == 0
}

// TaskSize checks TaskQueue size
func (c *Coordinator) TaskSize() int {
	return len(c.TaskQueue)
}

// Run runs in separate go thread and they are SEQUENTIAL
func (c *Coordinator) Run() {
	Wg.Add(1)
	for {
		select {
		case <-c.CTX.Done():
			Wg.Done()
			return
		default:
			if fun, data := c.Dequeue(); fun != nil {
				fun(data)
			}
		}
	}
}
