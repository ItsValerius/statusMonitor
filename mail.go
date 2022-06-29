package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

// SendMail sends an email to the given address using the given SMTP server.
func SendMail(address string, hostname string) {

	// Sender data.
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email address.
	to := []string{
		os.Getenv("MAIL_TO"),
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	fromStr := fmt.Sprintf("From: %s \r\n", from)
	toStr := fmt.Sprintf("To: %s \r\n", to[0])
	subject := fmt.Sprintf("Subject: %s Request failed \r\n\r\n", hostname)
	body := fmt.Sprintf("Request failed for: %s \r\nReported Time: %s \r\n", address, time.Now())
	msg := []byte(fromStr + toStr + subject + body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")
}
