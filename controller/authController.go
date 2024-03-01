package controller

import (
	_ "log"
	"strings"
	"time"

	"github.com/ADEMOLA200/danas-food/database"
	"github.com/ADEMOLA200/danas-food/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SingUp(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "invalid request",
		})
	}

	if user.Password != user.ConfirmPassword {
		c.Status(400)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "passwords do not match",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "error hashing password",
		})
	}
	user.Password = string(hashedPassword)

	r := database.DB.Create(&user)
	if r.Error != nil {
		if strings.Contains(r.Error.Error(), "Duplicate entry") {
			if strings.Contains(r.Error.Error(), "email") {
				c.Status(400)
				return c.JSON(fiber.Map{
					"success": false,
					"message": "error hashing password",
				})
			} else if strings.Contains(r.Error.Error(), "username") {
				c.Status(400)
				return c.JSON(fiber.Map{
					"success": false,
					"message": "username has been used",
				})
			}
		}

		c.Status(500)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "invalid request",
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "sucessful signup",
	})
}

func SignIn(c *fiber.Ctx) error {
	var loginRequest map[string]string

	if err := c.BodyParser(&loginRequest); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid JSON format.",
		})
	}

	var user models.User

	r := database.DB.Where("email = ?", loginRequest["password"]).First(&user)
	if r.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Wrong password!",
		})
	}

	c.Status(200)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "signIn successful",
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
