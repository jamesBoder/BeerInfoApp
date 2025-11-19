package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

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

// create a search command function
func searchCommand(s *state, cmd command) error {
	// base case : check if search term is provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("search term not provided")
	}

	// get the search term from command arguments
	var searchTerm string = cmd.args[0]

	// get the API key from config
	apiKey := ""
	if s != nil && s.config != nil {
		apiKey = s.config.apiKey
	}

	// exit if user types "quit" or "exit"
	if searchTerm == "quit" || searchTerm == "exit" {
		fmt.Println("Thanks for using Beer Info App! Goodbye!")
		return nil
	}

	// check if beer name is empty
	if searchTerm == "" {
		fmt.Println("Please enter a valid beer name.")
		return nil
	}

	// check if API key is set
	if apiKey == "" {
		fmt.Println("API key not set. Please set the API key in the configuration.")
		return nil
	}

	// API interaction section

	// build the URL with proper encoding
	baseURL := "https://beer9.p.rapidapi.com/"
	params := url.Values{}
	params.Add("name", searchTerm)
	fullURL := baseURL + "?" + params.Encode()

	// create a get request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("error creating request", err)
		return nil
	}

	// set the correct RapidAPI headers
	req.Header.Set("x-rapidapi-key", apiKey)
	req.Header.Set("x-rapidapi-host", "beer9.p.rapidapi.com")

	// make the request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request", err)
		return nil
	}

	//close the body of the response
	defer res.Body.Close()

	var apiResponse APIResponse // wrapper struct for large api responses

	// decode the body
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&apiResponse)
	if err != nil {
		fmt.Println("Error decoding parameters", err)
		return nil
	}

	// Check if the API returned an error
	if apiResponse.Error {
		fmt.Printf("API returned an error (code: %d)\n", apiResponse.Code)
		return nil
	}

	// Print the decoded beer information
	fmt.Printf("\n🍺 Found %d beers:\n\n", len(apiResponse.Data))
	for i, beer := range apiResponse.Data {
		fmt.Printf("--- Beer #%d ---\n", i+1)
		fmt.Printf("Name: %s\n", beer.Name)
		fmt.Printf("Brewery: %s\n", beer.Brewery)
		fmt.Printf("SKU: %s\n", beer.Sku)
		fmt.Printf("ABV: %s\n", beer.Abv)
		fmt.Printf("IBU: %s\n", beer.Ibu)
		fmt.Printf("Category: %s\n", beer.Category)
		fmt.Printf("Region: %s\n", beer.Region)
		fmt.Printf("Country: %s\n", beer.Country)
		fmt.Printf("Rating: %s\n", beer.Rating)
		fmt.Printf("Food Pairing: %s\n", beer.FoodPairing)
		fmt.Printf("Description: %s\n", beer.Description)
		fmt.Println()
	}

	// Ask if the user wants to search again
	fmt.Print("Would you like to search for another beer? (y/n): ")
	var again string
	fmt.Scanln(&again)
	if again == "no" || again == "n" || again == "N" || again == "No" || again == "NO" {
		fmt.Println("Thanks for using Beer Info App! Goodbye!")
		// exit the program
		os.Exit(0)
		return nil
	}

	fmt.Printf("Searching for beer %q using API key %s\n", searchTerm, apiKey)
	return nil

}

// create a help command function
func helpCommand(s *state, cmd command) error {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  login <username>   - Log in with the specified username")
	fmt.Println("  logout             - Log out of the current session")
	fmt.Println("  search <beername>  - Search for a specific beer")
	fmt.Println("  help               - Show this help message")
	fmt.Println("  exit, quit         - Exit the application")
	return nil
}

// create an exit command function
func exitCommand(s *state, cmd command) error {
	fmt.Println("Exiting the Beer Info App. Goodbye!")
	os.Exit(0)
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

// API Response wrapper structure
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
	Rating      string `json:"rating"`
	FoodPairing string `json:"food_pairing"`
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

	// greet the user
	fmt.Printf("Hello there, %s!\n", username)
	fmt.Println("*****************************")

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

	// create a config instance
	cfg := &config{
		apiKey: apiKey,
		User:   username,
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
	ch.registerCommand("search", searchCommand)
	ch.registerCommand("help", helpCommand)
	ch.registerCommand("exit", exitCommand)
	ch.registerCommand("quit", exitCommand)

	// CLI interaction section

	// Main Loop
	for {

		// display current user
		if s.config.User != "" {
			fmt.Printf("\nCurrent User: %s\n", s.config.User)
		} else {
			fmt.Println("\nCurrent User: guest")
		}

		// prompt user for command
		fmt.Println("\n> Enter a command (type 'help' for available commands): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// split input into command name and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue // skip empty input
		}

		// parse input into command struct
		cmd := command{
			name: parts[0],  // first part is command name
			args: parts[1:], // remaining parts are arguments
		}

		// run the command
		err := ch.runCommand(s, cmd)
		if err != nil {
			fmt.Println("Error executing command:", err)
			fmt.Println("Type 'help' to see the list of available commands.")
		}

		// end of main loop
		fmt.Println()

	}

}
