package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"net/smtp"

	"github.com/ADEMOLA200/danas-food/database"
	"github.com/ADEMOLA200/danas-food/models"
	"github.com/AfterShip/email-verifier"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	verifier = emailverifier.NewVerifier()
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
		"message": "Verification email will be sent to your email address",
	})
}

func SignIn(c *fiber.Ctx) error {
	var loginRequest map[string]string

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	var user models.User

	r := database.DB.Where("email = ?", loginRequest["email"]).First(&user)
	if r.Error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	// Generate and send OTP
	otp := generateOTP()
	err := sendOTPEmail(user.Email, otp)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error sending OTP"})
	}

	// Save OTP to the database (optional)
	saveOTPToDatabase(user.ID, otp)

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "OTP sent to your email"})
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
	// Set up SMTP authentication information
	APPLICATION_SPECIFIC_PASSWORD := "ddng hdkh odyq dapc"
	auth := smtp.PlainAuth("", "odukoyaabdullahi01@gmail.com", APPLICATION_SPECIFIC_PASSWORD, "smtp.gmail.com")

	// Compose the email message
	subject := "Your OTP for sign-in"
	body := fmt.Sprintf("Your OTP (One-Time Password) for sign-in is: %s", otp)
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Send the email
	err := smtp.SendMail("smtp.gmail.com:465", auth, "odukoyaabdullahi01@gmail.com", []string{email}, msg)
	if err != nil {
		return err
	}

	return nil
}

func saveOTPToDatabase(userID uint, otp string) {
	///////
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
