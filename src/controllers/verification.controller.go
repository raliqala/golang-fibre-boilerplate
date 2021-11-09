package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/database"
	"github.com/raliqala/golang-fibre-boilerplate/src/helpers"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
	"github.com/raliqala/golang-fibre-boilerplate/src/services"
)

func EmailVerification(c *fiber.Ctx) error {
	db := database.DB

	refreshToken := c.Params("token")

	response := services.ValidateToken(refreshToken)
	results := helpers.Unfold(response)

	if !results.Success {
		fmt.Println(results.Message)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": results.Message,
		})
	}

	if len(results.Data) == 0 {
		fmt.Println(results.Message)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": results.Message,
		})
	}

	var user models.User

	if userExists := db.Where("uuid = ?", results.Data).First(&user).RowsAffected; userExists <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Access Denied",
		})
	}

	if createErr := db.Model(&user).Where("uuid = ?", results.Data).Update("verified", true).Error; createErr != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Your account was successfully verified",
	})

}
