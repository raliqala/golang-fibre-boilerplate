package models

type User struct {
	Base
	FirstName string `json:"first_name" validate:"required,ascii,max=128"`
	LastName  string `json:"last_name" validate:"required,ascii,max=128"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,max=255"`
	Verified  bool   `json:"verified" gorm:"default:0"`
}
