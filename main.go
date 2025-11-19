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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	config             *config
	lastSearchResulsts []Beer // slice to hold last search results
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

	// save last search results to state
	s.lastSearchResulsts = apiResponse.Data

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

	// Ask user what to do next

	fmt.Println("\nOptions:")
	fmt.Println("  [s] Search again")
	fmt.Println("  [m] Main menu")
	fmt.Println("  [x] Exit")
	fmt.Print("Your choice: ")

	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	switch choice {
	case "s", "search":
		fmt.Println("Let's search for another beer!")
		// Loop continues
	case "m", "menu", "":
		fmt.Println("Returning to main menu...")
		// print help menu
		helpCommand(s, command{})
		return nil
	case "x", "exit", "quit":
		fmt.Println("Thanks for using Beer Info App! Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Returning to main menu...")
		return nil
	}
	return nil
}

// create a help command function
func helpCommand(s *state, cmd command) error {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  login <username>   - Log in with the specified username")
	fmt.Println("  logout             - Log out of the current session")
	fmt.Println("  search <beername>  - Search for a specific beer")
	fmt.Println("  favorite <beername> - Add a beer to your favorites")
	fmt.Println("  favorites          - Display your favorite beers")
	fmt.Println("  remove <beername>  - Remove a beer from your favorites")
	fmt.Println("  clear favs         - Clear all favorite beers")
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

// create a title case function for beer names
func toTitleCase(input string) string {
	caser := cases.Title(language.English)
	return caser.String(input)
}

// create a saveFavorites function to save favorite beers to a JSON file
func saveFavorites(favorites Favorites) error {
	// create or truncate the favorites file
	file, err := os.Create("favorites.json")
	if err != nil {
		return err
	}
	defer file.Close()

	// encode the favorites struct to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(favorites)
	if err != nil {
		return err
	}

	return nil
}

// create a loadFavorites function to load favorite beers from a JSON file
func loadFavorites() (Favorites, error) {
	var favorites Favorites

	// open the favorites file
	file, err := os.Open("favorites.json")
	if err != nil {
		return favorites, err
	}
	defer file.Close()

	// decode the JSON data into the favorites struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&favorites)
	if err != nil {
		return favorites, err
	}

	return favorites, nil
}

// create a favorite command function
func favoriteCommand(s *state, cmd command) error {
	// check if beer name is provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("beer name not provided")
	}

	// get the beer name from command arguments
	beerName := toTitleCase(strings.Join(cmd.args, " "))

	// load existing favorites
	favorites, err := loadFavorites()
	if err != nil {
		// if file not found, initialize empty favorites
		if os.IsNotExist(err) {
			favorites = Favorites{Beers: []Beer{}}
		} else {

			// other errors
			fmt.Println("Error loading favorites:", err)
			return nil
		}
	}

	// check if beer is already in favorites
	for _, beer := range favorites.Beers {
		if strings.EqualFold(beer.Name, beerName) {
			fmt.Printf("Beer %q is already in your favorites.\n", beerName)
			return nil
		}
	}

	// find the beer in last search results

	// initialize a pointer to hold the beer to add
	var beerToAdd *Beer
	found := false

	if len(s.lastSearchResulsts) > 0 {
		for _, beer := range s.lastSearchResulsts {
			if strings.EqualFold(beer.Name, beerName) {
				beerToAdd = &beer
				found = true
				break
			}
		}
	}

	if found {
		// save the beer to favorites
		favorites.Beers = append(favorites.Beers, *beerToAdd)
		fmt.Printf("Beer %q added to favorites with all details\n", beerName)
	} else {
		// save a beer with only the name if not found in last search results
		// create empty beer with only name
		newBeer := Beer{Name: beerName}
		// append to favorites
		favorites.Beers = append(favorites.Beers, newBeer)
		fmt.Printf("Beer %q added to favorites with name only\n", beerName)
	}

	// save updated favorites
	err = saveFavorites(favorites)
	if err != nil {
		fmt.Println("Error saving favorites:", err)
		return nil
	}

	fmt.Printf("Beer %q added to favorites!\n", beerName)
	return nil
}

// display favorite beers command function
func displayFavoritesCommand(s *state, cmd command) error {
	// load existing favorites
	favorites, err := loadFavorites()
	if err != nil {
		fmt.Println("error getting favorites", err)
		return nil
	}

	// check if there are any favorites
	if len(favorites.Beers) == 0 {
		fmt.Println("No favorite beers found.")
		return nil
	}

	// display favorite beers
	fmt.Println("\nYour Favorite Beers:")
	// iterate over favorite beers and print their names
	for i, beer := range favorites.Beers {
		fmt.Printf("\n--- Beer #%d ---\n", i+1)
		fmt.Printf("Name: %s\n", beer.Name)

		// Only show fields if they exist
		if beer.Brewery != "" {
			fmt.Printf("Brewery: %s\n", beer.Brewery)
		}
		if beer.Abv != "" {
			fmt.Printf("ABV: %s\n", beer.Abv)
		}
		if beer.Ibu != "" {
			fmt.Printf("IBU: %s\n", beer.Ibu)
		}
		if beer.Category != "" {
			fmt.Printf("Category: %s\n", beer.Category)
		}
		if beer.Region != "" {
			fmt.Printf("Region: %s\n", beer.Region)
		}
		if beer.Country != "" {
			fmt.Printf("Country: %s\n", beer.Country)
		}
		if beer.Description != "" {
			fmt.Printf("Description: %s\n", beer.Description)
		}
		fmt.Println()
	}

	return nil
}

// add a remove favorite command function
func removeFavoriteCommand(s *state, cmd command) error {
	// check if beer name is provided
	if len(cmd.args) == 0 {
		return fmt.Errorf("beer name not provided")
	}

	// get the beer name from command arguments
	beerName := toTitleCase(strings.Join(cmd.args, " "))

	// load existing favorites
	favorites, err := loadFavorites()
	if err != nil {
		fmt.Println("error loading favorites", err)
		return nil
	}

	// find and remove the beer from favorites (case-insensitive)
	index := -1
	for i, beer := range favorites.Beers {
		// compare case-insensitively using strings.EqualFold
		if strings.EqualFold(beer.Name, beerName) {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Printf("Beer %q not found in favorites.\n", beerName)
		return nil
	}

	// remove the beer from the slice
	favorites.Beers = append(favorites.Beers[:index], favorites.Beers[index+1:]...)

	// save updated favorites
	err = saveFavorites(favorites)
	if err != nil {
		fmt.Println("error saving favorites", err)
	}
	fmt.Printf("Beer %q removed from favorites!\n", beerName)
	return nil
}

// create a clear favorites command that removes all favorite beers
func clearFavoritesCommand(s *state, cmd command) error {
	// create an empty favorites struct
	favorites := Favorites{Beers: []Beer{}}

	// save the empty favorites to the file
	err := saveFavorites(favorites)
	if err != nil {
		fmt.Println("error clearing favorites", err)
		return nil
	}

	fmt.Println("All favorite beers have been cleared.")
	return nil
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

// create a favorites struct to hold favorite beers
type Favorites struct {
	Beers []Beer `json:"beers"`
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
	fmt.Printf("Hello there, %s!\n", toTitleCase(username))
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
	ch.registerCommand("favorite", favoriteCommand)
	ch.registerCommand("favorites", displayFavoritesCommand)
	ch.registerCommand("remove", removeFavoriteCommand)
	ch.registerCommand("clear favs", clearFavoritesCommand)

	// CLI interaction section

	// Main Loop
	for {

		// display current user
		if s.config.User != "" {
			fmt.Printf("\nCurrent User: %s\n", toTitleCase(s.config.User))
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
