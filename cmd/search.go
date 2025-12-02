package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jamesBoder/BeerInfoApp.git/internal/api"
	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// SearchCommand - implements cli.Command
type SearchCommand struct {
	ctx *CommandContext
}

// RandomCommand - implements cli.Command
type RandomCommand struct {
	ctx *CommandContext
}

// NewSearchCommand returns a instance of SearchCommand
func NewSearchCommand(ctx *CommandContext) *SearchCommand {
	return &SearchCommand{ctx: ctx}
}

// name method
func (c *SearchCommand) Name() string {
	return "search"
}

// Description method
func (c *SearchCommand) Description() string {
	return "Search by name"
}

// usage method
func (c *SearchCommand) Usage() string {
	return "search <beer name>"
}

// Execute method
func (c *SearchCommand) Execute(args []string) error {
	// base case : check if search term is provided
	if len(args) == 0 {
		fmt.Println(ui.Error("Error: Search term is required"))
		fmt.Println(ui.Tip("\nUsage: search <beer name or brewery>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   search IPA"))
		fmt.Println(ui.Info("   search Sixpoint"))
		fmt.Println(ui.Info("   search \"Hazy IPA\""))
		fmt.Println(ui.Tip("\nTip: Use quotes for multi-word searches"))
		return nil
	}

	// get the search term from command arguments
	var searchTerm string = args[0]

	// exit if user types "quit" or "exit"
	if searchTerm == "quit" || searchTerm == "exit" {
		ui.ShowGoodbyeMessage()
		return nil
	}

	// check if beer name is empty
	if searchTerm == "" {
		fmt.Print(ui.Prompt("Please enter a valid beer name."))
		return nil
	}

	// API interaction section

	// create a BeerAPI client
	client := api.NewBeerAPIClient(c.ctx.State.Config.APIKey)

	// call the SearchBeers method
	apiResponse, err := client.SearchBeers(searchTerm)
	if err != nil {
		fmt.Println(ui.Error("Error: Could not connect to beer database"))
		fmt.Println(ui.Info("Reason: %v", err))
		fmt.Println(ui.Info("\nPossible solutions:"))
		fmt.Println(ui.Info("   • Check your internet connection"))
		fmt.Println(ui.Info("   • Verify your API key is still valid"))
		fmt.Println(ui.Info("   • The API service might be temporarily down"))
		fmt.Println(ui.Info("   • Try again in a few moments"))
		fmt.Println(ui.Info("\n🌐 API Status: https://rapidapi.com/status"))
		return nil
	}

	// Check if the API returned an error
	if apiResponse.Error {
		fmt.Println(ui.Error("Error: Beer database returned an error"))
		fmt.Println(ui.Info("🔢 Error code: %d", apiResponse.Code))

		// switch on error code

		switch apiResponse.Code {
		case 400:
			fmt.Println(ui.Error("Error: Bad request"))
			fmt.Println(ui.Info("Reason: Invalid search parameters"))
		case 401:
			fmt.Println(ui.Error("Error: Unauthorized"))
			fmt.Println(ui.Info("Reason: Invalid API key or access denied"))
		case 403:
			fmt.Println(ui.Error("Error: Forbidden"))
			fmt.Println(ui.Info("Reason: Access to the resource is forbidden"))
		case 404:
			fmt.Println(ui.Error("Error: Not found"))
			fmt.Println(ui.Info("Reason: The requested resource was not found"))
		case 500:
			fmt.Println(ui.Error("Error: Internal server error"))
			fmt.Println(ui.Info("Reason: The server encountered an error processing your request"))
			fmt.Println(ui.Tip("   • This is not your fault"))
			fmt.Println(ui.Tip("   • Try again in a few minutes"))
		default:
			fmt.Println(ui.Error("Error: Unknown error occurred (code %d)", apiResponse.Code))
		}

		return nil
	}

	if len(apiResponse.Data) == 0 {
		fmt.Println(ui.Info("No beers found matching '%s'", searchTerm))
		fmt.Println(ui.Tip("\nSuggestions:"))
		fmt.Println(ui.Tip("   • Try a broader search term"))
		fmt.Println(ui.Tip("   • Check your spelling"))
		fmt.Println(ui.Tip("   • Try searching by brewery name"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   search IPA"))
		fmt.Println(ui.Info("   search Sixpoint"))
		return nil
	}

	// save last search results to models.State
	c.ctx.State.LastSearchResults = apiResponse.Data

	// save search history
	historyEntry := models.SearchHistoryEntry{
		Term:      searchTerm,
		Timestamp: time.Now(),
		Results:   apiResponse.Data,
	}
	c.ctx.State.SearchHistory = append(c.ctx.State.SearchHistory, historyEntry)

	// persist search history to file
	err = storage.SaveSearchHistory(c.ctx.State.Config.User, c.ctx.State.SearchHistory)
	if err != nil {
		fmt.Println(ui.Warning("Warning: Could not save search history: %v\n", err))
	}

	// Print the decoded beer information
	fmt.Println(ui.Header("\nFound %d beers:\n\n", len(apiResponse.Data)))
	for i, beer := range apiResponse.Data {
		fmt.Println(ui.FormatBeer(beer, i+1))
	}

	return nil
}

// NewRandomCommand returns a instance of RandomCommand
func NewRandomCommand(ctx *CommandContext) *RandomCommand {
	return &RandomCommand{ctx: ctx}
}

// name method
func (c *RandomCommand) Name() string {
	return "random"
}

// Description method
func (c *RandomCommand) Description() string {
	return "Get a random beer"
}

// usage method
func (c *RandomCommand) Usage() string {
	return "random"
}

// Random execute method
func (c *RandomCommand) Execute(args []string) error {
	// check if ther are any last search results
	if len(c.ctx.State.LastSearchResults) == 0 {
		fmt.Println(ui.Warning(" No last search results found. Please perform a search first."))
		fmt.Println(ui.Tip("\nTo use the random feature:"))
		fmt.Println(ui.Info("   1. First search for beers: search IPA"))
		fmt.Println(ui.Info("   2. Then get a random suggestion: random"))
		fmt.Println(ui.Tip("\nTip: The random command picks from your last search"))
		return nil
	}

	// seed the random number generator
	//rand.Seed(time.Now().UnixNano())

	// generate a random index
	randomIndex := rand.Intn(len(c.ctx.State.LastSearchResults))

	// get the random beer
	randomBeer := c.ctx.State.LastSearchResults[randomIndex]

	// display the random beer information
	fmt.Println(ui.Header("\nRandom Beer from Last Search Results:\n\n"))

	// add Formatted beer display
	fmt.Println(ui.FormatBeer(randomBeer, 1))

	return nil

}
