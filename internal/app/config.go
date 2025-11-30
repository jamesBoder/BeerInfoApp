package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Handle all config loading and validation

// create a config struct
type Config struct {
	APIKey   string
	Username string
}

// LoadConfig(), loads .env file, gets API key, validates API Key, and returns a config instance or error
func LoadConfig() (*Config, error) {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// get API key from .env file
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not found in .env file")
	}

	// create and return config instance
	config := &Config{
		APIKey: apiKey,
	}

	return config, nil
}

// PromptForUsername(), prompts user for username and returns it
func PromptForUsername() string {
	username := ""
	fmt.Println("Enter your username (or press Enter to continue as guest): ")
	fmt.Scanln(&username)

	return username
}
