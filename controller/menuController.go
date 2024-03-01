package controller

import (
    "github.com/gofiber/fiber/v2"
	"github.com/ADEMOLA200/danas-food/database"
    "github.com/ADEMOLA200/danas-food/models"
)

// AddMenuItem adds a new menu item
func AddMenuItem(c *fiber.Ctx) error {
    var menuItem models.MenuItem
    if err := c.BodyParser(&menuItem); err != nil {
        return err
    }
    database.DB.Create(&menuItem)
    return c.JSON(fiber.Map{"message": "Menu item added successfully", "menu_item": menuItem})
}

// UpdateMenuItem updates an existing menu item
func UpdateMenuItem(c *fiber.Ctx) error {
    id := c.Params("id")
    var updatedItem models.MenuItem
    if err := c.BodyParser(&updatedItem); err != nil {
        return err
    }
    var existingItem models.MenuItem
    result := database.DB.First(&existingItem, id)
    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Menu item not found"})
    }
    existingItem.Name = updatedItem.Name
    existingItem.Price = updatedItem.Price
    database.DB.Save(&existingItem)
    return c.JSON(existingItem)
}

// DeleteMenuItem deletes a menu item
func DeleteMenuItem(c *fiber.Ctx) error {
    id := c.Params("id")
    var menuItem models.MenuItem
    result := database.DB.First(&menuItem, id)
    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Menu item not found"})
    }
    database.DB.Delete(&menuItem)
    return c.SendStatus(fiber.StatusNoContent)
}
