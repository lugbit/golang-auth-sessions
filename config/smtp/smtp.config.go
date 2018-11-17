// Package smtp is used to set up the activation email sender email
package smtp

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail builds and sends the email using the sender and SMTP details above
func SendEmail(t string, rec string) {
	// Set up authentication information
	auth := smtp.PlainAuth("", os.Getenv("SENDER_ADDRESS"), os.Getenv("SENDER_PASSWORD"), os.Getenv("SMTP_SERVER"))
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	vLink := fmt.Sprintf("http://localhost:8080/activate-account?token=%v", t)
	to := []string{rec}
	msg := []byte("To:" + rec + "\r\n" +
		"Subject: Account Activation\r\n" +
		"\r\n" +
		"Please click the link below to activate your account:\n" + vLink + "\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("SENDER_ADDRESS"), to, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Activation email sent to ", rec)
}
