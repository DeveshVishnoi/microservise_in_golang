package utils

import (
	"errors"
	"log"
	"regexp"

	"github.com/joho/godotenv"
)

var EnvData map[string]string

func EnvConfig() {

	envFile, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error to load Env file: ", err)
	}
	EnvData = envFile
}

// ValidateEmail checks if the email format is valid
func ValidateEmail(email string) error {
	// Define the email regex pattern
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

	// Check if the email matches the pattern
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}
