package database

import (
	"fmt"
	"freecode/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicln("Error on .env file")
	}
	user := os.Getenv("USER_DB")
	password := os.Getenv("PASSWORD_DB")
	db_name := os.Getenv("DBNAME")

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s", user, password, db_name)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}
