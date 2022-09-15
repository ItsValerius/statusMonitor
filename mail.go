package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

// SendMail sends an email to the given address using the given SMTP server.
func SendMail(address string) {
	params := GetParams("parameters")
	// Sender data.
	from := MAIL_ADDR
	password := MAIL_PW
	customName := params["dclp"]

	// Receiver email address.
	//to := strings.Split(params["receiver"], ",")
	toHeader := params["receiver"]

	// smtp server configuration.
	smtpHost := params["smtp"]
	smtpPort := params["port"]
	fromStr := fmt.Sprintf("From: %s \r\n", from)
	toStr := fmt.Sprintf("To: %s \r\n", toHeader)
	subject := fmt.Sprintf("Subject: %s Service offline \r\n\r\n", customName)
	body := fmt.Sprintf("Request failed for: %s \r\nReported Time: %s \r\n", address, time.Now().Format(time.RFC3339))
	msg := []byte(fromStr + toStr + subject + body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		log.Panic(err)
	}
	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(toHeader); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
	log.Println("Email Sent Successfully!")

}
