package middleware

import (
	_"github.com/ADEMOLA200/danas-food/util"
	"github.com/gofiber/fiber/v2"
)

// Middleware to check if the user is authenticated
func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}