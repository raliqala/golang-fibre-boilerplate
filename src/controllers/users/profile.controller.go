package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/golang-fibre-boilerplate/src/database"
	"github.com/raliqala/golang-fibre-boilerplate/src/helpers"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
	"github.com/raliqala/golang-fibre-boilerplate/src/utils"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

func GetProfile(c *fiber.Ctx) error {
	db := database.DB
	id := c.Locals("id")

	u := new(models.User)
	if res := db.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Cannot find the User"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    u,
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	db := database.DB
	id := c.Locals("id")

	data := new(utils.UpdatePassword)

	if bodyErr := c.BodyParser(data); bodyErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   bodyErr,
		})
	}

	if ok, err := helpers.ValidateInput(*data); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	if err := passwordvalidator.Validate(data.NewPassword, 60); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"password": err.Error(),
			},
		})
	}

	var user models.User

	if userExists := db.Where("uuid = ?", id).First(&user).RowsAffected; userExists <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Incorrect credentials.",
		})
	}

	if ok := utils.CheckPasswordHash(data.OldPassword, user.Password); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect old password.",
		})
	}

	hash, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	data.NewPassword = string(hash)

	// check if is new password
	if check := utils.CheckPasswordHash(data.OldPassword, data.NewPassword); check {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Previously used password, please choose another",
		})
	}

	if createErr := db.Model(&user).Where("uuid = ?", id).Update("password", data.NewPassword).Error; createErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"general": "Your password updated successfully",
	})
}

//update email
func UpdateEmail() {

}
