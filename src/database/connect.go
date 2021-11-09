package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPg() {
	var err error

	p := config.Config("DB_PORT")
	port, err := strconv.Atoi(p)
	// port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		log.Println("Sorry db port error: ", err)
	}

	// connection url to DB
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))

	// connect to DB
	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database..")
	}

	fmt.Println("Running the migrations...")
	DB.AutoMigrate(&models.User{}, &models.Claims{})
	fmt.Println("Database connection was successful...")

}
