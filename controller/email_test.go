package controller

import (
    "errors"
    "fmt"

    "gopkg.in/gomail.v2"
)

// EmailSender
type EmailSender struct {
    Message *gomail.Message
    dailer  *gomail.Dialer
}

// NewEmailSender
func NewEmailSender() *EmailSender {
    m := gomail.NewMessage()
    dailer := gomail.NewDialer("smtp.gmail.com", 587, "smtpEmail@example.com", "smtp-password")
    m.SetHeader("From", "from-email@gmail.com")

    return &EmailSender{
        Message: m,
        dailer:  dailer,
    }
}

// Send
func (es *EmailSender) Send(subject, to, message string) error {

    es.Message.SetHeader("To", to)
    es.Message.SetHeader("Subject", subject)
    es.Message.SetBody("text/plain", message)

    if err := es.dailer.DialAndSend(es.Message); err != nil {
        return err
    }

    return nil
}

// checks before sending the email
//   - code is valid
//   - email already sent
//   - code is expired or not
//
// TODO : Implement your checks

func ChecksBeforeEmailSend(code string) error {
    return errors.New("email already sent")
}

func Mail() {
    code := "123456"

    emailer := NewEmailSender()

    if err := ChecksBeforeEmailSend(code); err != nil {
        panic(fmt.Sprintf("Check %v failed", err.Error()))
    }

    err := emailer.Send(
        "Verify your email", // subject
        "to@gmail.com", // to email
        fmt.Sprintf("Your verfication code is %v", code), // message
     )

    if err != nil {
        panic(err)
    }
}