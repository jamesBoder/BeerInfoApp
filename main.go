package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// get API key from .env file
	apiKey := os.Getenv("Api_Key")
	if apiKey == "" {
		fmt.Println("API key not set")
		return
	}

	// url set
	url := "https://rapidapi.com/"
	if url == "" {
		fmt.Println("No url provided")
		return
	}

	// create a get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error getting request", err)
		return
	}
}
