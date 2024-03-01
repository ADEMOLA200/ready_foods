package controller

import (
	_"crypto/rand"
	_"encoding/base32"
	"errors"
	"log"
	_"net/http"
	_"net/smtp"
	_"strings"
	"time"

	"github.com/ADEMOLA200/danas-food/database"
	_ "github.com/ADEMOLA200/danas-food/middleware"
	"github.com/ADEMOLA200/danas-food/models"
	"github.com/gofiber/fiber/v2"
	_ "github.com/sendgrid/sendgrid-go"
	_ "github.com/sendgrid/sendgrid-go/helpers/mail"
	_"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// func Register(c *fiber.Ctx) error {
// 	var user models.User
// 	if err := c.BodyParser(&user); err != nil {
// 		c.Status(400)
// 		return c.JSON(fiber.Map{
// 			"message": "Invalid request",
// 		})
// 	}

// 	if user.Password != user.ConfirmPassword {
// 		c.Status(400)
// 		c.JSON(fiber.Map{
// 			"message": "passwords are not the same",
// 		})
// 	}

// 	/// Generate OTP
// 	otp, err := generateOTP()
// 	if err != nil {
// 		c.Status(http.StatusInternalServerError)
// 		return c.JSON(fiber.Map{
// 			"message": "Failed to generate OTP",
// 		})
// 	}

// 	// Send OTP to user via email
// 	err = sendOTP(user.Email, otp)
// 	if err != nil {
// 		c.Status(http.StatusInternalServerError)
// 		return c.JSON(fiber.Map{
// 			"message": "Failed to send OTP",
// 		})
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
// 	if err != nil {
// 		c.Status(500)
// 		c.JSON(fiber.Map{
// 			"message": "Failed to register user",
// 		})
// 	}
// 	user.Password = string(hashedPassword)

// 	r := database.DB.Create(&user)
// 	if r.Error != nil {
// 		if strings.Contains(r.Error.Error(), "Duplicated entry") {
// 			if strings.Contains(r.Error.Error(), "username") {
// 				c.Status(400)
// 				return c.JSON(fiber.Map{
// 					"message": "username has been taken",
// 				})
// 			} else if strings.Contains(r.Error.Error(), "email") {
// 				c.Status(400)
// 				return c.JSON(fiber.Map{
// 					"message": "email has been used",
// 				})
// 			}
// 		}

// 		c.Status(500)
// 		return c.JSON(fiber.Map{
// 			"message": "Failed to register user",
// 		})
// 	}

// 	c.Status(200)
// 	return c.JSON(fiber.Map{
// 		"message": "Registered user successfully",
// 	})
// }

// // generateOTP generates a random 6-digit OTP
// func generateOTP() (string, error) {
// 	const otpLength = 6
// 	otpBytes := make([]byte, otpLength)
// 	_, err := rand.Read(otpBytes)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base32.StdEncoding.EncodeToString(otpBytes)[:otpLength], nil
// }

// // sendOTP sends the OTP to the user's email
// func sendOTP(email, otp string) error {
// 	// Connect to the SMTP server
// 	auth := smtp.PlainAuth("", "odukoyaabdullahi01@gmail.com", "xkyk lzbc muxr uptb", "smtp.gmail.com")

// 	// Compose the email message
// 	from := "odukoyaabdullah@gmail.com"
// 	to := email
// 	subject := "OTP Verification"
// 	body := "Your OTP: " + otp
// 	msg := []byte("To: " + to + "\r\n" +
// 		"From: " + from + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body + "\r\n")

// 	// Send the email
// 	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // ///////////////////////////////////////////////////////////////////////
// func TestSendOTP(c *fiber.Ctx) error {
// 	email := "odukoyaabdullahi01@gmail.com" // Replace with the recipient's email address
// 	otp := "123456"                         // Replace with the generated OTP
// 	err := sendOTP(email, otp)
// 	if err != nil {
// 		return err
// 	}
// 	return c.SendString("OTP sent successfully")
// }

// //////////////////////////////////////////////////////////////////////////////////

// func Login(c *fiber.Ctx) error {
// 	var loginRequest map[string]string
// 	if err := c.BodyParser(&loginRequest); err != nil {
// 		c.Status(400)
// 		c.JSON(fiber.Map{
// 			"message": "invalid request",
// 		})
// 	}

// 	var user models.User
// 	r := database.DB.Where("email = ?", loginRequest["email"]).First(&user)
// 	if r.Error != nil {
// 		c.Status(401)
// 		c.JSON(fiber.Map{
// 			"message": "Unauthorized user",
// 		})
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest["password"])); err != nil {
// 		c.Status(401)
// 		return c.JSON(fiber.Map{
// 			"message": "Unauthorized user",
// 		})
// 	}

// 	c.Status(200)
// 	return c.JSON(fiber.Map{
// 		"message": "Login user successfully",
// 	})
// }

// func Logout(c *fiber.Ctx) error {
// 	cookie := fiber.Cookie{
// 		Name:     "jwt",
// 		Value:    "",
// 		Expires:  time.Now().Add(-time.Hour),
// 		HTTPOnly: true,
// 	}

// 	c.Cookie(&cookie)

// 	return c.JSON(fiber.Map{
// 		"message": "Logout successful",
// 	})
// }

// Handlers for FoodOrder endpoints

func CreateOrder(c *fiber.Ctx) error {

	// Check if the user is authenticated
	// if err := middleware.IsAuthenticated(c); err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized user"})
	// }

	var newOrder models.FoodOrder
	if err := c.BodyParser(&newOrder); err != nil {
		return err
	}

	newOrder.CreateDtm = time.Now() // Set current time as creation time
	database.DB.Create(&newOrder)
	return c.JSON(newOrder)
}

func GetAllOrders(c *fiber.Ctx) error {

	// Check if the user is authenticated
	// if err := middleware.IsAuthenticated(c); err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized user"})
	// }

	var orders []*models.FoodOrder
	// Preload the associated menu items
	database.DB.Preload(clause.Associations).Find(&orders)
	return c.JSON(orders)
}

func GetOrder(c *fiber.Ctx) error {

	// Check if the user is authenticated
	// if err := middleware.IsAuthenticated(c); err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized user"})
	// }

	id := c.Params("id")
	var order models.FoodOrder
	result := database.DB.Preload("MenuItems").First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	}
	return c.JSON(order)
}

func UpdateOrder(c *fiber.Ctx) error {

	// Check if the user is authenticated
	// if err := middleware.IsAuthenticated(c); err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized user"})
	// }

	id := c.Params("id")

	// Parse the updated order from the request body
	var updatedOrder models.FoodOrder
	if err := c.BodyParser(&updatedOrder); err != nil {
		return err
	}

	// Retrieve the existing order from the database
	var existingOrder models.FoodOrder
	result := database.DB.Preload("MenuItems").First(&existingOrder, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
	}

	// Update order fields
	existingOrder.UserPhone = updatedOrder.UserPhone
	existingOrder.Name = updatedOrder.Name
	existingOrder.Address = updatedOrder.Address
	existingOrder.TotalItems = updatedOrder.TotalItems
	existingOrder.TotalPay = updatedOrder.TotalPay

	// Update menu items
	for i, updatedMenuItem := range updatedOrder.MenuItems {
		if i < len(existingOrder.MenuItems) {
			// Update existing menu item
			existingMenuItem := existingOrder.MenuItems[i]
			existingMenuItem.Name = updatedMenuItem.Name
			existingMenuItem.Price = updatedMenuItem.Price
			database.DB.Save(&existingMenuItem)
		} else {
			// Add new menu item
			database.DB.Create(&updatedMenuItem)
			existingOrder.MenuItems = append(existingOrder.MenuItems, updatedMenuItem)
		}
	}

	// Remove deleted menu items
	for i := len(existingOrder.MenuItems) - 1; i >= len(updatedOrder.MenuItems); i-- {
		existingMenuItem := existingOrder.MenuItems[i]
		database.DB.Delete(&existingMenuItem)
		existingOrder.MenuItems = append(existingOrder.MenuItems[:i], existingOrder.MenuItems[i+1:]...)
	}

	// Save the updated order
	if err := database.DB.Save(&existingOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update order"})
	}

	// Retrieve the updated order from the database with associated menu items
	updatedOrderWithMenuItems := models.FoodOrder{}
	database.DB.Preload("MenuItems").First(&updatedOrderWithMenuItems, id)

	return c.JSON(updatedOrderWithMenuItems)
}

func DeleteOrder(c *fiber.Ctx) error {

	// Check if the user is authenticated
	// if err := middleware.IsAuthenticated(c); err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized user"})
	// }

	id := c.Params("id")
	if id == "" {
		log.Println("Error: id parameter is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id parameter is required"})
	}

	var order models.FoodOrder
	result := database.DB.First(&order, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Error: Order with id %s not found\n", id)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Order not found"})
		} else {
			log.Printf("Error: %v\n", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
		}
	}

	// Delete menu items associated with the order
	var menuItems []models.MenuItem
	result = database.DB.Where("order_id = ?", id).Find(&menuItems)
	if result.Error != nil {
		log.Printf("Error: %v\n", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}
	for _, menuItem := range menuItems {
		result = database.DB.Delete(&menuItem)
		if result.Error != nil {
			log.Printf("Error: %v\n", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
		}
	}

	result = database.DB.Delete(&order)
	if result.Error != nil {
		log.Printf("Error: %v\n", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	log.Printf("Order with id %s has been deleted\n", id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order has been deleted successfully"})
}
