package mail

import (
	"log"

	"gopkg.in/gomail.v2"
)

const (
    SMTPHost = "localhost"
    SMTPPort = 1025
)

func Send(to, subject, body, attachmentPath string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "maxime@amazingmarket.com")
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    if attachmentPath != "" {
        m.Attach(attachmentPath)
    }

    d := gomail.NewDialer(SMTPHost, SMTPPort, "", "")

    err := d.DialAndSend(m)
    if err != nil {
        log.Printf("Error sending email: %v", err)
        return err
    }

    return nil
}
