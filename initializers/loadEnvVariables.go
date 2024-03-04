package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load(".env.local")
	if err != nil {
		env := os.Getenv("ENV")
		if env == "dev" {
			err := godotenv.Load(".env.dev")
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}
		if env == "prod" {
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
