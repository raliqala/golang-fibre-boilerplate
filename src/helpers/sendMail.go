package helpers

import (
	"crypto/tls"
	"log"
	"strconv"
	"time"

	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Payload struct {
	To            string `json:"to"`
	Name          string `json:"name"`
	Cc            string `json:"cc"`
	HTMLContent   string `json:"htmlContent"`
	TextContent   string `json:"textContent"`
	Subject       string `json:"subject"`
	URL           string `json:"url"`
	Inline        string `json:"inline"`
	FileName      string `json:"filename"`
	IsAttachment  bool   `json:"isAttachment"`
	IsTextContent bool   `json:"isTextContent"`
}

func SendEmail(user Payload) {
	server := mail.NewSMTPClient()

	port := config.Config("SMTP_PORT")
	smtp_port, err := strconv.Atoi(port)
	if err != nil {
		log.Println("Sorry db port error: ", err)
	}

	server.Host = config.Config("SMTP_HOST")
	server.Port = smtp_port
	server.Username = config.Config("SMTP_USER")
	server.Password = config.Config("SMTP_PASSWORD")
	server.Encryption = mail.EncryptionSTARTTLS

	server.KeepAlive = false
	server.ConnectTimeout = 20 * time.Second
	server.SendTimeout = 20 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom(config.Config("FROM_NAME") + "<" + config.Config("FROM_EMAIL") + ">").
		AddTo(user.Name + "<" + user.To + ">").
		AddCc(user.Cc).
		SetSubject(user.Subject)

	email.SetBody(mail.TextHTML, user.HTMLContent)

	if user.IsTextContent {
		email.AddAlternative(mail.TextPlain, user.TextContent)
	}

	if user.IsAttachment {
		email.Attach(&mail.File{FilePath: user.URL, Name: user.FileName, Inline: true})
	}

	if email.Error != nil {
		log.Fatal(email.Error)
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
}
