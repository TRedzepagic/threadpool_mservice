package ping

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	"github.com/TRedzepagic/threadpool_mservice/internal/mail"
	"github.com/TRedzepagic/threadpool_mservice/pkg/pool"
)

// Ping represents data contract for pinging
type Ping struct {
	IP           string   `json:"ip"`
	Recipients   []string `json:"recipients"`
	PingInterval string   `json:"pinginterval"`
}

// Func represents the pinging function
func Func(pingData []byte) {

	ping := Ping{}
	err := json.Unmarshal(pingData, &ping)
	if err != nil {
		log.Println("error unmarshalling ", err.Error())
	}
	log.Println(ping)

	// Ping syscall, -c is ping count, -i is interval, -w is timeout
	out, _ := exec.Command("ping", ping.IP, "-c 5", "-i "+ping.PingInterval, "-w 2").Output()

	if (strings.Contains(string(out), "Destination Host Unreachable")) || (strings.Contains(string(out), "100% packet loss")) {
		log.Printf("Host %s is down, sending mail ... \n", ping.IP)
		// Enqueue to send e-mail when able
		pool.CoordinatorInstance.Enqueue(mail.Func, pingData)
	} else {
		log.Println("Host ping successful!")
	}

	// Sleeping for debugging purposes, also to get some time to cancel if need be. Optional.
	// fmt.Println("Sleeping 3sec, ease of following timeline")
	// time.Sleep(3 * time.Second)

}
