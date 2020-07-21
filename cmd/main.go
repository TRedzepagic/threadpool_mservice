package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TRedzepagic/threadpool_mservice/pkg/pool"
)

type Mail struct {
	Address string `json:"address"`
}

type Ping struct {
	IP string `json:"ip_address"`
}

func PingFunc(data []byte) {
	ping := Ping{}
	json.Unmarshal(data, &ping)
	fmt.Println(ping.IP)

	// If ping fails
	mail := Mail{
		Address: "ph.ph2@gmail.com",
	}
	bytes, _ := json.Marshal(mail)
	Pool.Enqueue(MailFunc, bytes)
}

func MailFunc(data []byte) {
	mail := Mail{}
	json.Unmarshal(data, mail)
	fmt.Println("SENDING MAIL TO " + mail.Address)
}

var Pool pool.Coordinator

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	Pool.TaskQueue = make([]pool.Function, 0)
	Pool.DataQueue = make([][]byte, 0)
	Pool.Done = make(chan bool, 1)
	Pool.RunToMain = make(chan bool, 1)

	go Pool.Run()

	ping := Ping{IP: "127.0.0.1"}
	bytes, _ := json.Marshal(ping)
	Pool.Enqueue(PingFunc, bytes)

	<-stop
	Pool.Stop()
}
