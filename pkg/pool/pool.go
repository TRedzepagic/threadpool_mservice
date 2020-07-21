package pool

import (
	"fmt"
	"sync"
)

// This package contains implementation of Thread Pool that will be able like PoolCoordinator

// ByteArray -
type ByteArray [][]byte

// Function -
type Function func([]byte)

// Coordinator is implementation of Thread Pool that uses one queue for deploying and executing tasks
type Coordinator struct {
	TaskQueue []Function
	DataQueue ByteArray
	Done      chan bool
	RunToMain chan bool
	mux       sync.Mutex
}

// Enqueue places new task into queue and returns its lenghth
func (c *Coordinator) Enqueue(fun func([]byte), data []byte) int {
	c.mux.Lock()
	c.TaskQueue = append(c.TaskQueue, fun)
	c.DataQueue = append(c.DataQueue, data)
	c.mux.Unlock()
	return len(c.TaskQueue)
}

// Dequeue removes one task and retunrs it to caller
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

// IsEmpty checks if Coordinator queue is empty
func (c *Coordinator) IsEmpty() bool {
	return len(c.TaskQueue) == 0 || len(c.DataQueue) == 0
}

// Stop -
func (c *Coordinator) Stop() {
	c.Done <- true

	// This Stop is called from Main thread. We must not return from this function until
	// we receive signal on RunToMain Channel from Run thread. This signal will be received when Run is gracefully stoped.
	<-c.RunToMain
}

// Run runs in separate go thread
func (c *Coordinator) Run() {
	for {
		select {
		case <-c.Done:
			fmt.Println(" I am exiting ")
			c.RunToMain <- true
			break
		default:
			if !c.IsEmpty() {
				fun, data := c.Dequeue()
				fun(data)
			}
		}
	}
}
