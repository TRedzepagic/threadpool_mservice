package pool

import (
	"context"
	"sync"
)

// This package contains implementation of a Thread Pool.

// ByteSlice is a type redefinition for clarity - slice of byte slices ([]byte).
type ByteSlice [][]byte

// Function is a type redefinition for clarity - function that takes a byte slice ([]byte)
type Function func([]byte)

// Wg represents our WaitGroup variable, used for graceful exit.
// Waits for our worker to complete current task in queue, then exits
// Ignores non-started tasks in queue.
var Wg sync.WaitGroup

// Coordinator is implementation of Thread Pool that uses one queue for deploying and executing tasks
type Coordinator struct {
	TaskQueue []Function
	DataQueue ByteSlice
	Ctx       context.Context
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

// Size checks TaskQueue&DataQueue size
func (c *Coordinator) Size() (int, int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.TaskQueue), len(c.DataQueue)
}

// Run runs in separate go thread/worker, its subsequent tasks are SEQUENTIAL.
func (c *Coordinator) Run() {
	Wg.Add(1)
	for {
		select {
		case <-c.Ctx.Done():
			Wg.Done()
			return
		default:
			if fun, data := c.Dequeue(); fun != nil {
				fun(data)
			}
		}
	}
}
