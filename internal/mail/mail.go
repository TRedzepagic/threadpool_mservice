package mail

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		log.Println("error opening configuration", err.Error())
	}

	var mailconfiguration mailConfig
	err = json.Unmarshal(data, &mailconfiguration)
	if err != nil {
		log.Println("error unmarshalling ", err.Error())
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

	log.Println("Configuration of user mail: ")
	log.Println("Host: " + config.Host)
	log.Println("Mailing port: " + strconv.Itoa(config.Port))
	log.Println("User: Hidden")
	log.Println("Pass: Hidden")

	mail := mailer.NewMail()
	mail.FromName = "Go Mailer - Redzep Microservice"
	mail.From = config.User

	toMail := Mail{}
	err := json.Unmarshal(mailData, &toMail)
	if err != nil {
		log.Println("error unmarshalling ", err.Error())
	}

	// Debugging purposes. Gives user time to cancel the program, testing out the context cancellation and graceful exit.
	log.Println("Sleeping (10 seconds), try cancelling now!! Should finish running task..")
	time.Sleep(10 * time.Second)

	for _, recipientIteration := range toMail.Recipients {
		mail.SetTo(recipientIteration)
	}
	log.Println("Mailing function received the following:")
	log.Println(toMail)

	mail.Subject = "Admin notice : Server Down"
	mail.Body = "Your server is down. Host Address: " + toMail.IP + " " + "Host pinging interval:" + toMail.PingInterval

	log.Println("Not actually mailing, testing to avoid clutter : ")
	log.Println("Detected e-mails : ")
	log.Println(toMail.Recipients)

	// Used for actual mailing, uncomment when needed
	// mailerino := mailer.NewMailer(config, true)
	// err := mailerino.Send(mail)
	// if err != nil {
	// 	println(err.Error())
	// } else {
	// 	log.Println("Mail sent to : ")
	// 	log.Println(toMail.Recipients)
	// }
}
