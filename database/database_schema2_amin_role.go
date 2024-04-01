package database

import (
	"fmt"
	"knowledgeplus/go-api/initializers"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb02Admin() *gorm.DB {
	// init env
	initializers.LoadEnvVariables()

	// init database
	DB_USERNAME := os.Getenv("DB_USERNAME_ADMIN_ROLE")
	DB_PASSWORD := os.Getenv("DB_PASSWORD_ADMIN_ROLE")
	DB_NAME := os.Getenv("DB_NAME_02")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")

	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn: ", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return db
}
