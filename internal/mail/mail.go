package mail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/irnes/go-mailer"
)

// mailConfig struct represents sender configuration to be unmarshalled
type mailConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

// Mail data contract for function
type Mail struct {
	IP           string   `json:"ip"`
	Recipients   []string `json:"recipients"`
	PingInterval string   `json:"pinginterval"`
}

// initMailer initializes sender information and cofig
func initMailer(path string) *mailConfig {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error opening configuration", err.Error())
	}

	var mailconfiguration mailConfig
	err = json.Unmarshal(data, &mailconfiguration)

	if err != nil {
		fmt.Println("error unmarshalling ", err.Error())
	}
	return &mailconfiguration
}

// wrapperConf configuration struct necessary to avoid initialization repetition
type wrapperConf struct {
	mux        sync.Mutex
	configured bool
	config     mailer.Config
}

var wrapConf wrapperConf = wrapperConf{}

// Func represents the mailing function
func Func(mailData []byte) {
	wrapConf.mux.Lock()
	if !wrapConf.configured {
		fmt.Println("Mailconfig repetition printing test")
		mailConf := initMailer(os.Getenv("MAILCONFJSON"))
		wrapConf.config.Host = mailConf.Host
		wrapConf.config.Port = mailConf.Port
		wrapConf.config.User = mailConf.User
		wrapConf.config.Pass = mailConf.Pass
		wrapConf.configured = true
	}
	wrapConf.mux.Unlock()

	fmt.Println("Configuration of user mail: ")
	fmt.Println("Host: " + wrapConf.config.Host)
	fmt.Println("Mailing port: " + strconv.Itoa(wrapConf.config.Port))
	fmt.Println("User: Hidden")
	fmt.Println("Pass: Hidden")

	mail := mailer.NewMail()
	mail.FromName = "Go Mailer - Redzep Microservice"
	mail.From = wrapConf.config.User

	toMail := Mail{}
	json.Unmarshal(mailData, &toMail)

	// Debugging purposes. Gives user time to cancel the program, testing out the context cancellation and graceful exit.
	fmt.Println("Sleeping (10 seconds), try cancelling now!! Should finish running task..")
	time.Sleep(10 * time.Second)

	for _, recipientIteration := range toMail.Recipients {
		mail.SetTo(recipientIteration)
	}
	fmt.Println("Mailing function received the following:")
	fmt.Println(toMail)

	mail.Subject = "Admin notice : Server Down"
	mail.Body = "Your server is down. Host Address: " + toMail.IP + " " + "Host pinging interval:" + toMail.PingInterval

	fmt.Println("Not actually mailing, testing to avoid clutter : ")
	fmt.Println("Detected e-mails : ")
	fmt.Println(toMail.Recipients)

	// Used for actual mailing, uncomment when needed
	// mailerino := mailer.NewMailer(config, true)
	// err := mailerino.Send(mail)
	// if err != nil {
	// 	println(err.Error())
	// } else {
	// 	fmt.Println("Mail sent to : ")
	// 	fmt.Println(toMail.Recipients)
	// }
}
