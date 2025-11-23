package models

import (
	"time"
)

//---------- structs ----------//

// create config struct to hold app configuration
type Config struct {
	// API key string
	APIKey string
	// Username string
	User string
}

// create state struct that holds a pointer to a config
type State struct {
	// pointer to config struct
	Config            *Config
	LastSearchResults []Beer               // slice to hold last search results
	SearchHistory     []SearchHistoryEntry // slice to hold search history
}

// create a struct to hold search history entries
type SearchHistoryEntry struct {
	Term      string    // search term
	Timestamp time.Time // timestamp of the search
	Results   []Beer    // slice of beers returned in the search
}

// create command struct
type Command struct {
	Name string
	Args []string
}

// create a  struct to map command names to handler functions.
type CommandHandler struct {
	// map of command names to handler functions
	Handlers map[string]func(*State, Command) error
}

// ------------------ Structs for API Response ------------------ //

// API Response wrapper structure
type APIResponse struct {
	Code  int    `json:"code"`  // status code goes into apiResponse.Code
	Error bool   `json:"error"` // error field goes inot apiResponse.Code
	Data  []Beer `json:"data"`  // data array goes into apiResponse.Data which is a slice of beer structs
}

// Beer struct matching the actual API response
type Beer struct {
	Sku           string `json:"sku"`
	Name          string `json:"name"`
	Brewery       string `json:"brewery"`
	Description   string `json:"description"`
	Region        string `json:"region"`
	Country       string `json:"country"`
	Abv           string `json:"abv"`
	Ibu           string `json:"ibu"`
	Category      string `json:"category"`
	Rating        string `json:"rating"`
	FoodPairing   string `json:"food_pairing"`
	SubCategory_1 string `json:"sub_category_1"`
	SubCategory_2 string `json:"sub_category_2"`
}

// create a favorites struct to hold favorite beers
type Favorites struct {
	Beers []Beer `json:"beers"`
}

// RegisterCommand adds a command handler to the command handler map
func (ch *CommandHandler) RegisterCommand(name string, f func(*State, Command) error) {
	ch.Handlers[name] = f
}

// RunCommand executes a command by looking up its handler
func (ch *CommandHandler) RunCommand(s *State, cmd Command) error {
	// check if command exists in handlers map
	if handler, ok := ch.Handlers[cmd.Name]; ok {
		// call the handler function
		return handler(s, cmd)
	}
	return nil // command not found, but don't error
}
