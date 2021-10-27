package users

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/raliqala/safepass_api/src/config"
	"github.com/raliqala/safepass_api/src/database"
	"github.com/raliqala/safepass_api/src/helpers"
	"github.com/raliqala/safepass_api/src/models"
	"github.com/raliqala/safepass_api/src/services"
	"github.com/raliqala/safepass_api/src/utils"
)

func SignUp(c *fiber.Ctx) error {

	db := database.DB

	data := new(models.User)

	c.BodyParser(data)

	if ok, err := helpers.ValidateInput(*data); !ok {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	if userExists := db.Where(&models.User{Email: data.Email}).First(new(models.User)).RowsAffected; userExists > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"error":   "Sorry user already exists",
		})
	}

	hash, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	data.Password = string(hash)

	if createErr := db.Create(&data).Error; createErr != nil {
		return c.JSON(fiber.Map{
			"error":   true,
			"general": "Something went wrong, please try again later. 😕",
		})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := services.GenerateTokens(data.UUID.String())
	accessCookie, refreshCookie := services.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func SignIn(c *fiber.Ctx) error {
	db := database.DB

	data := new(utils.SignIn)

	c.BodyParser(data)

	if ok, err := helpers.ValidateInput(*data); !ok {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	user := new(models.User)

	if res := db.Where(&models.User{Email: user.Email}).First(user); res.RowsAffected <= 0 {
		c.JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect credentials",
		})
	}

	if ok := utils.CheckPasswordHash(data.Password, user.Password); !ok {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect credentials.",
		})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := services.GenerateTokens(user.UUID.String())
	accessCookie, refreshCookie := services.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

var jwtKey = []byte(config.Config("APP_SECRET"))

func GetAccessToken(c *fiber.Ctx) error {
	db := database.DB

	reToken := new(utils.RefreshToken)
	if err := c.BodyParser(reToken); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	refreshToken := reToken.RefreshToken

	refreshClaims := new(models.Claims)

	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if res := db.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ? AND ID = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer, refreshClaims.ID,
	).First(refreshClaims); res.RowsAffected <= 0 {
		// no such refresh token exist in the database
		c.ClearCookie("access_token", "refresh_token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			// refresh token is expired
			c.ClearCookie("access_token", "refresh_token")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
			})
		}
	} else {
		// malformed refresh token
		c.ClearCookie("access_token", "refresh_token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}

	_, accessToken := services.GenerateAccessClaims(refreshClaims.Issuer)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}
