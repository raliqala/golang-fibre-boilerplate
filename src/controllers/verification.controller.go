package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/database"
	"github.com/raliqala/golang-fibre-boilerplate/src/helpers"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
	"github.com/raliqala/golang-fibre-boilerplate/src/services"
	"github.com/raliqala/golang-fibre-boilerplate/src/utils"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func EmailVerification(c *fiber.Ctx) error {
	db := database.DB

	token := c.Params("token")

	response := services.ValidateToken(token)
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

func ForgotPassword(c *fiber.Ctx) error {
	db := database.DB

	data := new(utils.ForgotPass)

	c.BodyParser(data)

	if ok, err := helpers.ValidateInput(*data); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	var user models.User

	if userExists := db.Where("email = ?", data.Email).First(&user).RowsAffected; userExists <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "If this email exist, you will receive an email with instruction. Please follow them to reset your password",
		})
	}

	// setting up the email verification
	forgotPassword := services.GeneralTokens(user.UUID.String(), "forgot_password", 24)

	content := services.ResetPassword(utils.EmailVerification{
		Username:   user.FirstName,
		VerifyLink: forgotPassword,
	})

	helpers.SendEmail(helpers.Payload{
		To:          data.Email,
		Name:        user.FirstName,
		Cc:          "",
		HTMLContent: content,
		Subject:     "Reset your password",
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"error":   "If this email exist, you will receive an email with instruction. Please follow them to reset your password",
	})

}

func ResetPassword(c *fiber.Ctx) error {
	db := database.DB

	token := c.Params("token")

	data := new(utils.ResetPass)

	c.BodyParser(data)

	response := services.ValidateToken(token)

	results := helpers.Unfold(response)

	if !results.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": results.Message,
		})
	}

	if len(results.Data) == 0 {
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

	if err := passwordvalidator.Validate(data.Password, 60); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"password": err.Error(),
			},
		})
	}

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	data.Password = string(hash)

	if createErr := db.Model(&user).Where("uuid = ?", results.Data).Update("password", data.Password).Error; createErr != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Your password was updated successfully",
	})
}
