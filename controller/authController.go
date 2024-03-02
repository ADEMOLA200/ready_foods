package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/ADEMOLA200/danas-food/database"
	"github.com/ADEMOLA200/danas-food/models"
	"github.com/AfterShip/email-verifier"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	verifier = emailverifier.NewVerifier()
)

const (
	smtpUser     = "odukoyaabdullahi01@gmail.com"
	smtpPassword = "CB6DC2CC0E675DA892EC58BE0DC8D29BD301"
	smtpHost     = "smtp.elasticemail.com"
	smtpPort     = 2525
)

func SignUp(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if user.Password != user.ConfirmPassword {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	result, err := verifier.Verify(user.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error verifying email",
		})
	}

	if !result.Syntax.Valid {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Invalid email syntax",
		})
	}

	if result.Disposable {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Disposable email not allowed",
		})
	}

	if result.Reachable == "no" {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Email address not reachable",
		})
	}

	if !result.HasMxRecords {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Domain not properly set up to receive emails",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error hashing password",
		})
	}
	user.Password = string(hashedPassword)

	r := database.DB.Create(&user)
	if r.Error != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error creating user",
		})
	}

	// Send verification email
	err = sendVerificationEmail(user.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error sending verification email",
		})
	}

	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Verification email sent successfully",
	})
}

func SignIn(c *fiber.Ctx) error {
	var loginRequest map[string]string

	if err := c.BodyParser(&loginRequest); err != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	var user models.User

	r := database.DB.Where("email = ?", loginRequest["email"]).First(&user)
	if r.Error != nil {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Generate and send OTP
	otp := generateOTP()
	err := sendOTPEmail(user.Email, otp)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error sending OTP",
		})
	}

	// Send verification email
	err = sendVerificationEmail(user.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error sending verification email",
		})
	}

	// Save OTP to the database (optional)
	// saveOTPToDatabase(user.ID, otp)

	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"message": "OTP sent to your email",
	})
}

func generateOTP() string {
	// Generate a random 6-digit OTP
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(rune('0' + rand.Intn(10)))
	}
	return otp
}

func sendOTPEmail(email, otp string) error {
	// Use Elastic Email's SMTP server and port
	server := "smtp.elasticemail.com"
	port := 2525

	// Set up SMTP password
	PASSWORD := "CB6DC2CC0E675DA892EC58BE0DC8D29BD301"

	// Connect to the SMTP server without TLS encryption
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		return err
	}
	defer client.Close()

	// Set up SMTP authentication information after connecting
	auth := smtp.PlainAuth("", "odukoyaabdullahi01@gmail.com", PASSWORD, server)

	// Start TLS encryption
	if err := client.StartTLS(nil); err != nil {
		return err
	}

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return err
	}

	// Compose the email message
	subject := "Your OTP for sign-in"
	body := fmt.Sprintf("Your OTP (One-Time Password) for sign-in is: %s", otp)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Send the email
	if err := client.Mail("odukoyaabdullahi01@gmail.com"); err != nil {
		return err
	}
	if err := client.Rcpt(email); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func sendVerificationEmail(email string) error {

	// SMTP Auth
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Connect to SMTP Server
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		return err
	}

	// Enable TLS encryption
	client.StartTLS(nil)

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set sender and recipient
	if err = client.Mail(smtpUser); err != nil {
		return err
	}

	if err = client.Rcpt(email); err != nil {
		return err
	}

	// Send verification email
	msg := []byte("Subject: Verify Email\r\n" +
		"To: " + email + "\r\n\r\n" +
		"Please click to verify your email: http://example.com/verify?email=" + email)

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()

}

// func sendVerificationEmail(email string) error {
// 	// Similar to sendOTPEmail but for sending verification email
// 	// Use Elastic Email's SMTP server and port
// 	server := "smtp.elasticemail.com"
// 	port := 2525

// 	// Set up SMTP authentication information
// 	API_KEY := "70DF5D7505B200C6469A36F379D73DC465AEFAF0090697D8EA6AE2DACCB4BAC6BEED005114D216CB8423837FE5FA6135"

// 	auth := smtp.PlainAuth("", "odukoyaabdullahi01@gmail.com", API_KEY, server)

// 	// Compose the email message for verification
// 	subject := "Verify your email address"
// 	body := "Please verify your email address by clicking the link below:\n\n" +
// 		"http://odukoyaabdullahi01@gmail.com/verify?email=" + email
// 	msg := []byte("To: " + email + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body)

// 	// Send the email
// 	err := smtp.SendMail(fmt.Sprintf("%s:%d", server, port), auth, "odukoyaabdullahi01@gmail.com", []string{email}, msg)
// 	if err != nil {
// 		fmt.Printf("Error sending OTP email: %v\n", err)
// 		return err
// 	}

// 	return nil
// }

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
