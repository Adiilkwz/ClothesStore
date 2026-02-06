package mailer

import (
	"fmt"
	"log"
	"net/smtp"

	"clothes-store/internal/config"
)

var EmailQueue = make(chan string, 100)

func StartEmailWorker(cfg *config.Config) {
	go func() {
		for toEmail := range EmailQueue {
			sendEmail(cfg, toEmail)
		}
	}()
}

func sendEmail(cfg *config.Config, to string) {
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, cfg.SMTPHost)

	subject := "Order Confirmation"
	body := "Thank you for your order! It is being processed."

	msg := []byte(fmt.Sprintf("To: %s\r\n"+"Subject: %s\r\n"+"\r\n"+"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	err := smtp.SendMail(addr, auth, cfg.SMTPEmail, []string{to}, msg)

	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
	} else {
		log.Printf("Email sent succesfully to %s", to)
	}
}
