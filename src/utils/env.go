package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvironmentVariables() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	} else {
		log.Println(".env loaded")
	}
}
