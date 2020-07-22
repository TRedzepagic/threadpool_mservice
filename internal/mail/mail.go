package mail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/irnes/go-mailer"
)

type mailConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

// Mail data contract
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
	configured := false
	var config mailer.Config
	if !configured {
		mailConf := initMailer(os.Getenv("MAILCONFJSON"))
		config.Host = mailConf.Host
		config.Port = mailConf.Port
		config.User = mailConf.User
		config.Pass = mailConf.Pass
		configured = true
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

	for _, recipientIteration := range toMail.Recipients {
		mail.SetTo(recipientIteration)
	}
	fmt.Println("Mailing function received the following:")
	fmt.Println(toMail)

	mail.Subject = "Admin notice : Server Down"
	mail.Body = "Your server is down. Host Address: " + toMail.IP + " " + "Host pinging interval:" + toMail.PingInterval

	// fmt.Println("Not actually mailing, testing to avoid clutter : ")
	// fmt.Println("Detected e-mails : ")
	// fmt.Println(toMail.Recipients)

	// used for actual mailing, uncomment when needed

	mailerino := mailer.NewMailer(config, true)
	err := mailerino.Send(mail)
	if err != nil {
		println(err.Error())
	} else {
		fmt.Println("Mail sent to : ")
		fmt.Println(toMail.Recipients)
	}

}
