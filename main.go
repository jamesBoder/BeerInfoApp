package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// build command functions here

// create config struct to hold app configuration
type config struct {
	// API key string
	apiKey string
	// Username string
	User string
}

// create state struct that holds a pointer to a config
type state struct {
	// pointer to config struct
	config *config
}

// create command struct
type command struct {
	name string
	args []string
}

// create a  struct to map command names to handler functions.
type commandHandler struct {
	// map of command names to handler functions
	handlers map[string]func(*state, command) error
}

// create a login command functino
func loginCommand(s *state, cmd command) error {
	// check if username argument is provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}

	// get the username from the command arguments
	username := cmd.args[0]
	// set the username in the config
	s.config.User = username
	// print a success message
	fmt.Printf("User %s logged in successfully\n", username)
	return nil
}

// create a logout command function
func logoutCommand(s *state, cmd command) error {
	// check if user is logged in
	if s.config.User == "" {
		return fmt.Errorf("no user is currently logged in")
	}

	// get the current username
	username := s.config.User
	// clear the username in the config
	s.config.User = ""
	// print a success message
	fmt.Printf("User %s logged out successfully\n", username)
	return nil
}

// create a register handler function attached to commandHandler struct
func (ch *commandHandler) registerCommand(name string, f func(*state, command) error) {
	ch.handlers[name] = f

}

// create a run command function attached to commandHandler struct
func (ch *commandHandler) runCommand(s *state, cmd command) error {
	// check if command exists in handlers map
	if handler, ok := ch.handlers[cmd.name]; ok {
		// call the handler function
		return handler(s, cmd)

	}
	return fmt.Errorf("command not found")
}

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

	fmt.Println("Welcome to the Beer Info App!")
	fmt.Println("-----------------------------")

	// prompt the user to login or continue as guest

	username := ""
	fmt.Print("Enter your username (or press Enter to continue as guest): ")
	fmt.Scanln(&username)

	// If username is empty, set to "Guest"
	if username == "" {
		username = "guest"
	}
	fmt.Printf("Hello there, %s!\n", username)

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

	// CLI interaction section

	// Main Loop
	for {

		// prompt the user to enter a beer name
		var beerName string
		fmt.Print("Enter beer name (or 'quit' to exit) ")
		fmt.Scanln(&beerName)

		// exit if user types "quit" or "exit"
		if beerName == "quit" || beerName == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// check if beer name is empty
		if beerName == "" {
			fmt.Println("Please enter a valid beer name.")
			continue
		}

		// API interaction section

		// build the URL with proper encoding
		baseURL := "https://beer9.p.rapidapi.com/"
		params := url.Values{}
		params.Add("name", beerName)
		fullURL := baseURL + "?" + params.Encode()
		url := fullURL

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

		// Ask if the user wants to search again
		fmt.Print("Would you like to search for another beer? (y/n): ")
		var again string
		fmt.Scanln(&again)
		if again == "no" || again == "n" || again == "N" || again == "No" || again == "NO" {
			fmt.Println("Thanks for using Beer Info App! Goodbye!")
			break
		}

		// end of main loop
		fmt.Println()

	}

	// create a config instance
	cfg := &config{
		apiKey: apiKey,
		User:   "",
	}

	// create a state instance
	s := &state{
		config: cfg,
	}

	// create a commandHandler instance
	ch := &commandHandler{
		handlers: make(map[string]func(*state, command) error),
	}

	// register commands
	ch.registerCommand("login", loginCommand)
	ch.registerCommand("logout", logoutCommand)

	// Example command execution
	cmd := command{
		name: "login",
		args: []string{username},
	}

	err = ch.runCommand(s, cmd)
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

}
