package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/raliqala/golang-fibre-boilerplate/src/config"
	"github.com/raliqala/golang-fibre-boilerplate/src/models"
	"github.com/raliqala/golang-fibre-boilerplate/src/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPg() {
	p := config.Config("DB_PORT")
	port, err := strconv.Atoi(p)
	// port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Println("Sorry db port error: ", err)
	}

	// connection url to DB
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))

	// connect to DB
	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if dbErr != nil {
		panic("failed to connect to database..")
	}

	fmt.Println("Running the migrations...")

	if migrateErr := DB.AutoMigrate(&models.User{}, &models.Claims{}); migrateErr != nil {
		fmt.Println("Sorry couldn't migrate'...")
	}

	fmt.Println("Database connection was successful...")
}

// MI : An instance of MongoInstance Struct
var MI utils.MongoInstance

// ConnectDB - database connection
func ConnectMDB() {
	url := fmt.Sprintf("%s%s?%s", config.Config("MONGO_URI"), config.Config("DATABASE_NAME"), config.Config("MONGO_PREVILAGE"))

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connection was successful...")

	MI = utils.MongoInstance{
		Client: client,
		DB:     client.Database(config.Config("DATABASE_NAME")),
	}
}

func DBConnection(databaseName string) {
	switch databaseName {
	case "mongodb":
		ConnectMDB()
	case "postgres":
		ConnectPg()
	default:
		fmt.Println("No database name specified..")
	}
}
