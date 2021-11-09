// SecureAuth returns a middleware which secures all the private routes
package middleware

import (
	"time"

	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(config.Config("APP_SECRET"))

func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		claims := new(models.Claims)

		token, err := jwt.ParseWithClaims(accessToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"general": "Token Expired",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// this is not even a token, we should delete the cookies here
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return c.SendStatus(fiber.StatusUnauthorized)
			} else {
				// cannot handle this token
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			}
		}

		c.Locals("id", claims.Issuer)
		return c.Next()
	}
}
