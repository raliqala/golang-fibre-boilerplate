package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	"github.com/raliqala/golang-fibre-boilerplate/src/database"
	"github.com/raliqala/golang-fibre-boilerplate/src/helpers"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
	"github.com/raliqala/golang-fibre-boilerplate/src/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(config.Config("APP_SECRET"))

// GenerateTokens returns the access and refresh tokens
func GenerateTokens(uuid string) (string, string, int64, int64) {
	claim, accessToken, accessTime := GenerateAccessClaims(uuid)
	refreshToken, refreshTime := GenerateRefreshClaims(claim)

	return accessToken, refreshToken, accessTime, refreshTime
}

func GenerateAccessClaims(uuid string) (*models.Claims, string, int64) {

	t := time.Now()
	accessTime := t.Add(15 * time.Minute).Unix()

	claim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid,
			ExpiresAt: accessTime,
			Subject:   "access_token",
			IssuedAt:  t.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return claim, tokenString, accessTime
}

// GenerateRefreshClaims returns refresh_token
func GenerateRefreshClaims(cl *models.Claims) (string, int64) {

	t := time.Now()

	refreshTime := t.Add(30 * 24 * time.Hour).Unix()

	refreshClaim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: refreshTime,
			Subject:   "refresh_token",
			IssuedAt:  t.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)

	// create a claim on DB
	SaveToken(cl.Issuer, "refresh_token", refreshTime, refreshTokenString)

	if err != nil {
		panic(err)
	}

	return refreshTokenString, refreshTime
}

// SecureAuth returns a middleware which secures all the private routes
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

// GetAuthCookies sends two cookies of type access_token and refresh_token
func GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {
	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	return accessCookie, refreshCookie
}

// **************TODO new slate for token management below**************

// generate token
func GenerateToken(uuid string, subjectType string, timeExp int64) string {

	t := time.Now()
	claim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid,
			ExpiresAt: timeExp,
			Subject:   subjectType,
			IssuedAt:  t.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func SaveToken(uuid string, subjectType string, timeExp int64, token string) bool {
	db := database.DB

	t := time.Now()
	dataClaim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid,
			ExpiresAt: timeExp,
			Subject:   subjectType,
			IssuedAt:  t.Unix(),
			Audience:  token,
		},
	}

	// create a claim on DB
	if err := db.Create(&dataClaim).Error; err != nil {
		return false
	}

	return true
}

// For internal verification
func ValidateToken(accessToken string) string {

	db := database.DB

	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims := token.Claims.(*jwt.StandardClaims)

	if token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return helpers.ErrorHandle()
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			// this is not even a token, we should delete the cookies here
			return helpers.ErrorHandle()
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return helpers.ErrorHandle()
		} else {
			// cannot handle this token
			return helpers.ErrorHandle()
		}
	}

	if tokenExists := db.Where(&models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   claims.Issuer,
			Audience: accessToken,
			Subject:  claims.Subject,
		},
	}).First(&models.Claims{}); tokenExists.RowsAffected <= 0 {
		// no such refresh token exist in the database
		return helpers.ErrorHandle()
	}

	SuccessArray := utils.Success{Success: true, Message: "", Data: claims.Issuer}

	byteArray, err := json.Marshal(SuccessArray)
	if err != nil {
		fmt.Println(err)
	}

	if tokenExists := db.Where(&models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   claims.Issuer,
			Audience: accessToken,
			Subject:  claims.Subject,
		},
	}).Delete(&models.Claims{}); tokenExists.RowsAffected <= 0 {
		// no such refresh token exist in the database
		return helpers.ErrorHandle()
	}

	return string(byteArray)

}

// func VerifyToken(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt")

// 	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})

// 	if err != nil {
// 		c.Status(fiber.StatusUnauthorized)
// 		return c.JSON(fiber.Map{
// 			"message": "unauthenticated",
// 		})
// 	}

// 	claims := token.Claims.(*jwt.StandardClaims)

// if token.Valid {
// 	if claims.ExpiresAt < time.Now().Unix() {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error":   true,
// 			"general": "Token Expired",
// 		})
// 	}
// } else if ve, ok := err.(*jwt.ValidationError); ok {
// 	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
// 		// this is not even a token, we should delete the cookies here
// 		c.ClearCookie("access_token", "refresh_token")
// 		return c.SendStatus(fiber.StatusForbidden)
// 	} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
// 		// Token is either expired or not active yet
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	} else {
// 		// cannot handle this token
// 		c.ClearCookie("access_token", "refresh_token")
// 		return c.SendStatus(fiber.StatusForbidden)
// 	}
// }

// 	var user models.User

// 	database.DB.Where("id = ?", claims.Issuer).First(&user)

// 	return c.JSON(user)
// }

func GeneralTokens(uuid string, subject string, timeMultiplier time.Duration) string {

	t := time.Now()

	generalTime := t.Add(timeMultiplier * time.Hour).Unix()

	generalToken := GenerateToken(uuid, subject, generalTime)
	SaveToken(uuid, subject, generalTime, generalToken)

	return generalToken
}
