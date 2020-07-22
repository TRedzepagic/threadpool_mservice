package ping

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

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

	pingObject := Ping{}
	json.Unmarshal(pingData, &pingObject)
	fmt.Println("PingFunc received the following:")
	fmt.Println(pingObject)

	// Ping syscall, -c is ping count, -i is interval, -w is timeout
	out, _ := exec.Command("ping", pingObject.IP, "-c 5", "-i "+pingObject.PingInterval, "-w 2").Output()
	//fmt.Println("PING INTERVAL OF : " + pingObject.IP + " SET TO : " + pingObject.PingInterval)

	if (strings.Contains(string(out), "Destination Host Unreachable")) || (strings.Contains(string(out), "100% packet loss")) {
		fmt.Println("HOST " + pingObject.IP + " IS DOWN, SENDING MAIL..")
		// pinging failed, instantiate a structure representing host info, shorthand
		pool.CoordinatorInstance.Enqueue(mail.Func, pingData)

	} else {
		fmt.Println("Host ping successful! Ignoring mailing protocol.")
		fmt.Println(string(out))
	}
	fmt.Println("Sleeping 3sec, ease of following timeline")
	time.Sleep(3 * time.Second)

}
