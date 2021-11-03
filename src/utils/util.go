package utils

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadTemplates(template_path string) string {
	var data []byte
	switch template_path {
	case "email_validation":
		{
			file, err := os.ReadFile("/tmp/dat")
			check(err)
			data = file
		}
	}

	return string(data)

}
