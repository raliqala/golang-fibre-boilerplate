package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey"`
}
