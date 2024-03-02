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
	"crypto/tls"
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
	smtpPort     = 2525
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

	// Check if email is verified
	// if !user.EmailVerified {
	// 	c.Status(http.StatusBadRequest)
	// 	return c.JSON(fiber.Map{
	// 		"error": "Email not verified",
	// 	})
	// }

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
    // Set up TLS configuration
    tlsConfig := &tls.Config{
        ServerName: smtpHost,
    }

    // Connect to the SMTP server with TLS
    client, err := smtp.Dial(smtpHost + ":" + strconv.Itoa(smtpPort))
    if err != nil {
        return fmt.Errorf("failed to connect to SMTP server: %v", err)
    }
    defer client.Close()

    // Start TLS encryption
    if err := client.StartTLS(tlsConfig); err != nil {
        return fmt.Errorf("failed to start TLS: %v", err)
    }

    // Authentication
    auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
    if err := client.Auth(auth); err != nil {
        return fmt.Errorf("authentication failed: %v", err)
    }

    // Compose the email message
    subject := "Your OTP for sign-in"
    body := fmt.Sprintf("Your OTP (One-Time Password) for sign-in ready_food is: %s", otp)
    msg := []byte("From: Your Sender Name <mailtrap@demomailtrap.com>\r\n" + // Change this to your sender address
        "To: " + email + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "\r\n" +
        body)

    // Send the email
    if err := client.Mail("mailtrap@demomailtrap.com"); err != nil { // Change this to your sender address
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
