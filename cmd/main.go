package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TRedzepagic/threadpool_mservice/pkg/pool"
)

// Mail represents data contract for mailing
type Mail struct {
	Address string `json:"address"`
}

// Ping represents data contract for pinging
type Ping struct {
	IP string `json:"ip_address"`
}

// PingFunc represents the pinging function
func PingFunc(data []byte) {
	ping := Ping{}
	json.Unmarshal(data, &ping)
	fmt.Println(ping.IP)

	// If ping fails
	mail := Mail{
		Address: "ph.ph2@gmail.com",
	}
	bytes, _ := json.Marshal(mail)
	coordinator.Enqueue(MailFunc, bytes)
}

// MailFunc represents the mailing function
func MailFunc(data []byte) {
	mail := Mail{}
	json.Unmarshal(data, &mail)
	fmt.Println("SENDING MAIL TO " + mail.Address)
}

var coordinator = pool.CreateCoordinator()

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	context, stopCoordinator := context.WithCancel(context.Background())
	coordinator.CTX = context

	go coordinator.Run()

	ping := Ping{IP: "127.0.0.1"}
	bytes, _ := json.Marshal(ping)
	coordinator.Enqueue(PingFunc, bytes)

	<-stop
	stopCoordinator()
}
