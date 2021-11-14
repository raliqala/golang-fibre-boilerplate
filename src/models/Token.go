package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Claims struct {
	jwt.StandardClaims
	ID     uint   `gorm:"primaryKey"`
	Sector string `json:"sector"`
}

func (claims *Claims) BeforeCreate(tx *gorm.DB) error {
	// uuid.New() creates a new random UUID or panics.
	claims.Sector = "default"

	return nil
}
