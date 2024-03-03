package controller

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"text/template"
	"time"

	"crypto/tls"

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
	smtpUser     = "example@gmail.com" // Add your SMTP username here
    smtpPassword = "*********************" // Add your SMTP password here
    smtpHost     = "smtp.gmail.com" // Add your SMTP host here i'm smtp.gmail.com for this example
    smtpPort     = 2525 // Add your SMTP port here, you can get it from your provider
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

// Controller logic for sending OTP after successful sign-in attempt
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

    // Verify user's password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest["password"]))
    if err != nil {
        c.Status(http.StatusUnauthorized)
        return c.JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    // Generate OTP
    otp, err := generateOTP()
    if err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to generate OTP",
        })
    }

    // Log the generated OTP
    fmt.Println("Generated OTP:", otp)

    // Associate OTP with user
    user.OTP = otp
    if err := database.DB.Save(&user).Error; err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to save OTP",
        })
    }

    // Send OTP to user's email
    err = sendOTPEmail(user.Email, otp)
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

// Modified generateOTP function to return an error
func generateOTP() (string, error) {
	// Generate a random 6-digit OTP
	otp := ""
	for i := 0; i < 6; i++ {
		otp += string(rune('0' + rand.Intn(10)))
	}
	if len(otp) != 6 {
		return "", fmt.Errorf("failed to generate OTP")
	}
	return otp, nil
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

    // Load the HTML email template
    htmlTemplate := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>OTP Email</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f4f4f4;
					padding: 20px;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: #fff;
					padding: 20px;
					border-radius: 10px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				.header {
					background-color: #3498db;
					color: #fff;
					text-align: center;
					padding: 10px 0;
					border-top-left-radius: 10px;
					border-top-right-radius: 10px;
				}
				.content {
					padding: 20px;
				}
				.footer {
					text-align: center;
					padding: 10px 0;
					border-bottom-left-radius: 10px;
					border-bottom-right-radius: 10px;
				}
				.otp {
					font-size: 24px;
					text-align: center;
					margin-bottom: 20px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h2> Ready  Foods </h2>
				</div>
				<div class="content">
					<p>Dear User,</p>
					<p>Your OTP (One-Time Password) for sign-in is:</p>
					<div class="otp">{{.OTP}}</div>
					<p>Please use this OTP to proceed with your sign-in.</p>
				</div>
				<div class="footer">
					<p>This is an automated email. Please do not reply.</p>
				</div>
			</div>
		</body>
		</html>
    `

    // Create a new template and parse the HTML
    t := template.New("emailTemplate")
    t, err = t.Parse(htmlTemplate)
    if err != nil {
        return fmt.Errorf("failed to parse email template: %v", err)
    }

    // Prepare data to be passed into the template
    data := struct {
        OTP string
    }{
        OTP: otp,
    }

    // Execute the template to generate the HTML body
    var tpl bytes.Buffer
    if err := t.Execute(&tpl, data); err != nil {
        return fmt.Errorf("failed to execute template: %v", err)
    }
    htmlBody := tpl.String()

    // Compose the email message
    fromAddress := "Odukoya Abdullahi Ademola <no-reply@example.com>"
    toAddress := email
    subject := "Your OTP for sign-in"
    contentType := "text/html; charset=UTF-8"
    msg := []byte("From: " + fromAddress + "\r\n" +
        "To: " + toAddress + "\r\n" +
        "Subject: " + subject + "\r\n" +
        "MIME-Version: 1.0\r\n" +
        "Content-Type: " + contentType + "\r\n" +
        "\r\n" +
        htmlBody)

    // Send the email
    if err := client.Mail(smtpUser); err != nil {
        return fmt.Errorf("failed to send MAIL command: %v", err)
    }
    if err := client.Rcpt(toAddress); err != nil {
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

// Controller logic for verifying OTP
func VerifyOTP(c *fiber.Ctx) error {
    var verifyRequest map[string]string

    if err := c.BodyParser(&verifyRequest); err != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid JSON format",
        })
    }

    var user models.User

    // Fetch the user based on the OTP
    r := database.DB.Where("otp = ?", verifyRequest["otp"]).First(&user)
    if r.Error != nil {
        c.Status(http.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "Invalid OTP",
        })
    }

    // Clear OTP after successful verification
    user.OTP = ""
    if err := database.DB.Save(&user).Error; err != nil {
        c.Status(http.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "Failed to clear OTP",
        })
    }

    // Proceed with sign-in
    // You can set up session or JWT token here

    c.Status(http.StatusOK)
    return c.JSON(fiber.Map{
        "message": "OTP verification successful. You can now sign in.",
    })
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