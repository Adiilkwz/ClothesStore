package mailer

import (
	"fmt"
	"log"
	"net/smtp"

	"clothes-store/internal/config"
)

type EmailJob struct {
	To      string
	Subject string
	Body    string
}

var EmailQueue = make(chan EmailJob, 100)

func StartEmailWorker(cfg *config.Config) {
	go func() {
		for job := range EmailQueue {
			sendEmail(cfg, job)
		}
	}()
}

func sendEmail(cfg *config.Config, job EmailJob) {
	auth := smtp.PlainAuth("", cfg.SMTPEmail, cfg.SMTPPassword, cfg.SMTPHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", job.To, job.Subject, job.Body))

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	err := smtp.SendMail(addr, auth, cfg.SMTPEmail, []string{job.To}, msg)

	if err != nil {
		log.Printf("Failed to send email to %s: %v", job.To, err)
	} else {
		log.Printf("Email sent successfully to %s: %s", job.To, job.Subject)
	}
}
