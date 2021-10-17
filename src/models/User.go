package models

import "database/sql"

type User struct {
	Base
	Email       string `json:"email" gorm:"unique"`
	Username    string `json:"username" gorm:"unique"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Password    string `json:"password"`
	ActivatedAt sql.NullTime
}
