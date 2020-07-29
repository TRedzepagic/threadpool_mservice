package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/TRedzepagic/threadpool_mservice/internal/ping"
	"github.com/TRedzepagic/threadpool_mservice/pkg/pool"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	context, stopCoordinator := context.WithCancel(context.Background())
	pool.CoordinatorInstance.Ctx = context

	// Adds workers equal to the number of CPUs
	for i := 0; i < runtime.GOMAXPROCS(runtime.NumCPU()); i++ {
		go pool.CoordinatorInstance.Run()
	}

	// One worker
	// go pool.CoordinatorInstance.Run()

	pingInfo := ping.Ping{
		// Google IP address guaranteed to pass
		IP:           "172.217.16.100",
		Recipients:   make([]string, 0),
		PingInterval: "3"}
	// Add recipients
	pingInfo.Recipients = append(pingInfo.Recipients, "redzepagict@gmail.com")

	pingInfoSecond := ping.Ping{
		IP:           "192.168.0.25",
		Recipients:   make([]string, 0),
		PingInterval: "3"}
	// Add recipients
	pingInfoSecond.Recipients = append(pingInfoSecond.Recipients, "redzepagict@gmail.com")

	pingInfoThird := ping.Ping{
		IP:           "192.168.0.25",
		Recipients:   make([]string, 0),
		PingInterval: "3"}
	// Add recipients
	pingInfoThird.Recipients = append(pingInfoThird.Recipients, "redzepagict@gmail.com")

	pingBytes, _ := json.Marshal(pingInfo)
	pingBytesSecond, _ := json.Marshal(pingInfoSecond)
	pingBytesThird, _ := json.Marshal(pingInfoThird)

	pool.CoordinatorInstance.Enqueue(ping.Func, pingBytes)
	pool.CoordinatorInstance.Enqueue(ping.Func, pingBytesSecond)
	pool.CoordinatorInstance.Enqueue(ping.Func, pingBytesThird)

	<-stop
	stopCoordinator()
	pool.Wg.Wait()
	fmt.Println("Work done, shutting down")
}
