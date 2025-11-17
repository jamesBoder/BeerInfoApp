package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// API Response structure
type APIResponse struct {
	Code  int    `json:"code"`  // status code goes into apiResponse.Code
	Error bool   `json:"error"` // error field goes inot apiResponse.Code
	Data  []Beer `json:"data"`  // data array goes into apiResponse.Data which is a slice of beer structs
}

// Beer struct matching the actual API response
type Beer struct {
	Sku         string `json:"sku"`
	Name        string `json:"name"`
	Brewery     string `json:"brewery"`
	Description string `json:"description"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	Abv         string `json:"abv"`
	Ibu         string `json:"ibu"`
	Category    string `json:"category"`
}

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// get API key from .env file
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API key not set")
		return
	}

	// Set the correct Beer9 API URL
	url := "https://beer9.p.rapidapi.com/?brewery=Berkshire%20brewing%20company"

	// create a get request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating request", err)
		return
	}

	// set the correct RapidAPI headers
	req.Header.Set("x-rapidapi-key", apiKey)
	req.Header.Set("x-rapidapi-host", "beer9.p.rapidapi.com")

	// make the request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request", err)
		return
	}

	//close the body of the response
	defer res.Body.Close()

	var apiResponse APIResponse // wrapper struct for large api responses

	// decode the body
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResponse)
	if err != nil {
		fmt.Println("Error decoding parameters", err)
		return
	}

	// Check if the API returned an error
	if apiResponse.Error {
		fmt.Printf("API returned an error (code: %d)\n", apiResponse.Code)
		return
	}

	// Print the decoded beer information
	fmt.Printf("\n🍺 Found %d beers from Berkshire Brewing Company:\n\n", len(apiResponse.Data))
	for i, beer := range apiResponse.Data {
		fmt.Printf("--- Beer #%d ---\n", i+1)
		fmt.Printf("Name: %s\n", beer.Name)
		fmt.Printf("SKU: %s\n", beer.Sku)
		fmt.Printf("ABV: %s\n", beer.Abv)
		fmt.Printf("IBU: %s\n", beer.Ibu)
		fmt.Printf("Category: %s\n", beer.Category)
		fmt.Printf("Region: %s\n", beer.Region)
		fmt.Printf("Country: %s\n", beer.Country)
		fmt.Printf("Description: %s\n", beer.Description)
		fmt.Println()
	}
}
