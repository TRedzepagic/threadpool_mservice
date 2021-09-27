package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/TRedzepagic/threadpool_mservice/internal/ping"
	"github.com/TRedzepagic/threadpool_mservice/pkg/pool"
)

// load loads from config file
func load(path string) *ping.Pings {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("error opening configuration", err.Error())
	}

	var hosts ping.Pings
	err = json.Unmarshal(data, &hosts)
	if err != nil {
		log.Println("error unmarshalling ", err.Error())
	}
	return &hosts
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	context, stopCoordinator := context.WithCancel(context.Background())
	pool.CoordinatorInstance.Ctx = context

	HostSlice := load("config.json")

	// Adds workers equal to the number of CPUs
	for i := 0; i < runtime.GOMAXPROCS(runtime.NumCPU()); i++ {
		go pool.CoordinatorInstance.Run()
	}

	// One worker
	// go pool.CoordinatorInstance.Run()

	for _, OneHostInfo := range HostSlice.Hosts {
		pingBytes, err := json.Marshal(OneHostInfo)
		if err != nil {
			log.Println("error unmarshalling ", err.Error())
		}
		pool.CoordinatorInstance.Enqueue(ping.Func, pingBytes)
	}

	<-stop
	stopCoordinator()
	pool.Wg.Wait()
	log.Println("Work done, shutting down")
}
