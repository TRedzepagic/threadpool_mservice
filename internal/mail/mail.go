package mail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

// Func represents the mailing function
func Func(mailData []byte) {

	mailConf := initMailer(os.Getenv("MAILCONFJSON"))
	config := mailer.Config{
		Host: mailConf.Host,
		Port: mailConf.Port,
		User: mailConf.User,
		Pass: mailConf.Pass,
	}

	fmt.Println("Configuration of user mail: ")
	fmt.Println("Host: " + config.Host)
	fmt.Println("Mailing port: " + strconv.Itoa(config.Port))
	fmt.Println("User: Hidden")
	fmt.Println("Pass: Hidden")

	mail := mailer.NewMail()
	mail.FromName = "Go Mailer - Redzep Microservice"
	mail.From = config.User

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
