package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.StandardClaims
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}
