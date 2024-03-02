// package controller

// import (
// 	"fmt"
// 	"net/smtp"
// 	"testing"
// )

// func TestSendOTPEmail(t *testing.T) {
// 	email := "odukoyaabdullahi01@gmail.com" // Replace with your Gmail address
// 	otp := "123456"                  // Replace with a sample OTP

// 	err := sendOTPEmail(email, otp)
// 	if err != nil {
// 		fmt.Println("Error sending OTP:", err)
// 	} else {
// 		fmt.Println("OTP sent successfully")
// 	}
// }

// func sendOTPEmail(email, otp string) error {
// 	smtpUser := "odukoyaabdullahi01@gmail.com"
// 	smtpPassword := "*******************"
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := 587

// 	// Connect to the SMTP server with TLS encryption
// 	client, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
// 	if err != nil {
// 		return err
// 	}
// 	defer client.Close()

// 	// Set up SMTP authentication information after connecting
// 	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

// 	// Start TLS encryption
// 	if err := client.StartTLS(nil); err != nil {
// 		return err
// 	}

// 	// Authenticate
// 	if err := client.Auth(auth); err != nil {
// 		return err
// 	}

// 	// Compose the email message
// 	subject := "Your OTP for sign-in"
// 	body := fmt.Sprintf("Your OTP (One-Time Password) for sign-in is: %s", otp)
// 	msg := []byte("To: " + email + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body)

// 	// Send the email
// 	if err := client.Mail(smtpUser); err != nil {
// 		return fmt.Errorf("error sending MAIL command: %s", err)
// 	}
// 	if err := client.Rcpt(email); err != nil {
// 		return fmt.Errorf("error sending RCPT command: %s", err)
// 	}
// 	w, err := client.Data()
// 	if err != nil {
// 		return fmt.Errorf("error sending DATA command: %s", err)
// 	}
// 	_, err = w.Write(msg)
// 	if err != nil {
// 		return fmt.Errorf("error writing email body: %s", err)
// 	}
// 	err = w.Close()
// 	if err != nil {
// 		return fmt.Errorf("error closing email writer: %s", err)
// 	}

// 	return nil
// }

package controller

import "fmt"

func test() {
	fmt.Print("test")
}




// package controller

// import (
//     "errors"
//     "fmt"

//     "gopkg.in/gomail.v2"
// )

// // EmailSender
// type EmailSender struct {
//     Message *gomail.Message
//     dailer  *gomail.Dialer
// }

// // NewEmailSender
// func NewEmailSender() *EmailSender {
//     m := gomail.NewMessage()
//     dailer := gomail.NewDialer("smtp.gmail.com", 587, "smtpEmail@example.com", "smtp-password")
//     m.SetHeader("From", "from-email@gmail.com")

//     return &EmailSender{
//         Message: m,
//         dailer:  dailer,
//     }
// }

// // Send
// func (es *EmailSender) Send(subject, to, message string) error {

//     es.Message.SetHeader("To", to)
//     es.Message.SetHeader("Subject", subject)
//     es.Message.SetBody("text/plain", message)

//     if err := es.dailer.DialAndSend(es.Message); err != nil {
//         return err
//     }

//     return nil
// }

// // checks before sending the email
// //   - code is valid
// //   - email already sent
// //   - code is expired or not
// //
// // TODO : Implement your checks

// func ChecksBeforeEmailSend(code string) error {
//     return errors.New("email already sent")
// }

// func main() {
//     code := "123456"

//     emailer := NewEmailSender()

//     if err := ChecksBeforeEmailSend(code); err != nil {
//         panic(fmt.Sprintf("Check %v failed", err.Error()))
//     }

//     err := emailer.Send(
//         "Verify your email", // subject
//         "to@gmail.com", // to email
//         fmt.Sprintf("Your verfication code is %d", code), // message
//      )

//     if err != nil {
//         panic(err)
//     }
// }