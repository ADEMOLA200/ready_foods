package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
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
	smtpUser     = "api"
	smtpPassword = "a82e753a256d5c200074ddd37941735c"
	smtpHost     = "live.smtp.mailtrap.io"
	smtpPort     = 587
	authentication =	"plain"
	enable_starttls_auto = true
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

	c.Status(http.StatusOK)
	return c.JSON(fiber.Map{
		"message": "User created successfully",
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
	fmt.Println("Generated OTP:", otp)
	err := sendOTPEmail(user.Email, otp)
	if err != nil {
		fmt.Println("Error sending OTP email:", err)
		c.Status(http.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "Error sending OTP",
		})
	}

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
	// Add debug log
    fmt.Println("Sending OTP email to:", email)

	// Connect to the SMTP server
	client, err := smtp.Dial(smtpHost + ":" + strconv.Itoa(smtpPort))
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer client.Close()

	// Start TLS encryption
	if err := client.StartTLS(nil); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Compose the email message
	subject := "Your OTP for sign-in"
	body := fmt.Sprintf("Your OTP (One-Time Password) for sign-in is: %s", otp)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Send the email
	if err := client.Mail(smtpUser); err != nil {
		return fmt.Errorf("failed to send MAIL command: %v", err)
	}
	if err := client.Rcpt(email); err != nil {
		return fmt.Errorf("failed to send RCPT command: %v", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %v", err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	return nil
}

// func sendVerificationEmail(email string) error {
// 	// Use Mailtrap's SMTP server and port
// 	server := smtpHost
// 	port := smtpPort

// 	// Connect to the SMTP server without TLS encryption
// 	client, err := smtp.Dial(fmt.Sprintf("%s:%d", server, port))
// 	if err != nil {
// 		return err
// 	}
// 	defer client.Close()

// 	// Set up SMTP authentication information after connecting
// 	auth := smtp.PlainAuth("", smtpUser, smtpPassword, server)

// 	// Start TLS encryption
// 	if err := client.StartTLS(nil); err != nil {
// 		return err
// 	}

// 	// Authenticate
// 	if err := client.Auth(auth); err != nil {
// 		return err
// 	}

// 	// Compose the email message
// 	msg := []byte("Subject: Verify Email\r\n" +
// 		"To: " + email + "\r\n\r\n" +
// 		"Please click to verify your email: http://example.com/verify?email=" + email)

// 	// Send the email
// 	if err := client.Mail(smtpUser); err != nil {
// 		return err
// 	}
// 	if err := client.Rcpt(email); err != nil {
// 		return err
// 	}
// 	w, err := client.Data()
// 	if err != nil {
// 		return err
// 	}
// 	_, err = w.Write(msg)
// 	if err != nil {
// 		return err
// 	}
// 	err = w.Close()
// 	if err != nil {
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
