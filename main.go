package main

import (
	"bufio"
	"fmt"
	"math/rand"

	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/jamesBoder/BeerInfoApp.git/internal/api"
	"github.com/jamesBoder/BeerInfoApp.git/internal/export"
	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// build command functions here

// create a login command functino
func loginCommand(s *models.State, cmd models.Command) error {
	// check if username argument is provided
	if len(cmd.Args) == 0 {
		fmt.Println(ui.Error("Error: Username is required"))
		fmt.Println(ui.Tip("\nUsage: login <username>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   login james"))
		fmt.Println(ui.Info("   login sarah"))
		fmt.Println(ui.Tip("\nTip: Press Enter at startup to continue as guest"))
		return nil
	}

	// get the username from the command arguments
	username := cmd.Args[0]
	// set the username in the config
	s.Config.User = username

	// load search history
	history, err := storage.LoadSearchHistory(username)
	if err != nil {
		// if file not found, initialize empty history
		if os.IsNotExist(err) {
			s.SearchHistory = []models.SearchHistoryEntry{}
		} else {
			// other errors
			fmt.Println(ui.Error("error loading search history", err))
			return nil
		}
	} else {
		s.SearchHistory = history
	}

	// print a success message
	fmt.Println(ui.Success("User %s logged in successfully\n", username))
	return nil
}

// create a logout command function
func logoutCommand(s *models.State, cmd models.Command) error {
	// check if user is logged in
	if s.Config.User == "" {
		fmt.Println(ui.Warning("No user is currently logged in."))
		fmt.Println(ui.Tip("\n To log in, use the command: login <username>"))
		return nil
	}

	// get the current username
	username := s.Config.User

	// clear the username in the config
	s.Config.User = ""

	//clear last search results
	s.LastSearchResults = []models.Beer{}

	// clear search history
	s.SearchHistory = []models.SearchHistoryEntry{}

	// save empty favorites for guest user
	err := storage.SaveFavorites("guest", models.Favorites{Beers: []models.Beer{}})
	if err != nil {
		fmt.Println(ui.Error("error saving guest favorites", err))
		return nil
	}

	// save empty search history for guest user
	err = storage.SaveSearchHistory("guest", []models.SearchHistoryEntry{})
	if err != nil {
		fmt.Println(ui.Error("error saving guest search history", err))
		return nil
	}

	// print a success message
	fmt.Println(ui.Success("User %s logged out successfully\n", username))
	fmt.Println(ui.Success("Your favorite beers are saved and will be available when you log back in."))
	return nil
}

// create a search command function
func searchCommand(s *models.State, cmd models.Command) error {
	// base case : check if search term is provided
	if len(cmd.Args) == 0 {
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
	var searchTerm string = cmd.Args[0]

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
	client := api.NewBeerAPIClient(s.Config.APIKey)

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
	s.LastSearchResults = apiResponse.Data

	// save search history
	historyEntry := models.SearchHistoryEntry{
		Term:      searchTerm,
		Timestamp: time.Now(),
		Results:   apiResponse.Data,
	}
	s.SearchHistory = append(s.SearchHistory, historyEntry)

	// persist search history to file
	err = storage.SaveSearchHistory(s.Config.User, s.SearchHistory)
	if err != nil {
		fmt.Println(ui.Warning("Warning: Could not save search history: %v\n", err))
	}

	// Print the decoded beer information
	fmt.Println(ui.Header("\nFound %d beers:\n\n", len(apiResponse.Data)))
	for i, beer := range apiResponse.Data {
		fmt.Println(ui.FormatBeer(beer, i+1))
	}

	// Ask user what to do next

	fmt.Println(ui.Header("\nOptions:"))
	fmt.Println(ui.Info("  [s] Search again"))
	fmt.Println(ui.Info("  [m] Main menu"))
	fmt.Println(ui.Info("  [x] Exit"))
	fmt.Print(ui.Prompt("Your choice: "))
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
		helpCommand(s, models.Command{})
		return nil
	case "x", "exit", "quit":
		// exit the app
		ui.ShowGoodbyeMessage()
	default:
		fmt.Println("Returning to main menu...")
		return nil
	}
	return nil
}

// create a help command function
func helpCommand(s *models.State, cmd models.Command) error {
	// use FormatHelpMenu to display help menu
	fmt.Println(ui.FormatHelpMenu())
	return nil
}

// create an exit command function
func exitCommand(s *models.State, cmd models.Command) error {
	ui.ShowGoodbyeMessage()

	return nil
}

// create a favorite command function
func favoriteCommand(s *models.State, cmd models.Command) error {
	// check if beer name is provided
	if len(cmd.Args) == 0 {
		fmt.Println(ui.Error("Error: Beer name is required"))
		fmt.Println(ui.Tip("\nUsage: favorite <beer name>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   favorite Sixpoint Resin"))
		fmt.Println(ui.Info("   favorite \"Hazy IPA\""))
		fmt.Println(ui.Tip("\nTip: The beer should be from your last search results"))
		fmt.Println(ui.Info("   Or type 'search <name>' first to find a beer"))
		return nil
	}

	// get the beer name from command arguments
	beerName := ui.ToTitleCase(strings.Join(cmd.Args, " "))

	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		// if file not found, initialize empty favorites
		if os.IsNotExist(err) {
			favorites = models.Favorites{Beers: []models.Beer{}}
		} else {

			// other errors
			fmt.Println(ui.Error("Error loading favorites:", err))
			return nil
		}
	}

	// check if beer is already in favorites
	for _, beer := range favorites.Beers {
		if strings.EqualFold(beer.Name, beerName) {
			fmt.Println(ui.Warning("Beer %q is already in your favorites.\n", beerName))
			return nil
		}
	}

	// find the beer in last search results

	// initialize a pointer to hold the beer to add
	var beerToAdd *models.Beer
	found := false

	if len(s.LastSearchResults) > 0 {
		for _, beer := range s.LastSearchResults {
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
		fmt.Println(ui.Success("Beer %q added to favorites with all details\n", beerName))
	} else {
		// save a beer with only the name if not found in last search results
		// create empty beer with only name
		newBeer := models.Beer{Name: beerName}
		// append to favorites
		favorites.Beers = append(favorites.Beers, newBeer)
		fmt.Println(ui.Success("Beer %q added to favorites with name only\n", beerName))
	}

	// save updated favorites
	err = storage.SaveFavorites(s.Config.User, favorites)
	if err != nil {
		fmt.Println(ui.Error("Error saving favorites:", err))
		return nil
	}

	fmt.Println(ui.Success("Beer %q added to favorites!\n", beerName))
	return nil
}

// display favorite beers command function
func displayFavoritesCommand(s *models.State, cmd models.Command) error {
	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		fmt.Println(ui.Warning("favorites is empty. Type 'help' to add a favorite beer", err))
		return nil
	}

	// check if there are any favorites
	if err != nil || len(favorites.Beers) == 0 {
		fmt.Println(ui.Warning("You don't have any favorite beers yet"))
		fmt.Println(ui.Tip("\nHow to add favorites:"))
		fmt.Println(ui.Info("   1. Search for a beer: search IPA"))
		fmt.Println(ui.Info("   2. Add to favorites: favorite <beer name>"))
		fmt.Println(ui.Tip("\nOr type 'help' to see all commands"))
		return nil
	}

	// display favorite beers
	fmt.Println(ui.Header("\nYour Favorite Beers:"))
	// iterate over favorite beers and print their names
	for i, beer := range favorites.Beers {
		fmt.Println(ui.FormatBeerSummary(beer, i+1))
	}

	return nil
}

// add a remove favorite command function
func removeFavoriteCommand(s *models.State, cmd models.Command) error {
	// check if beer name is provided
	if len(cmd.Args) == 0 {
		fmt.Println(ui.Warning("beer name not provided"))
		fmt.Println(ui.Tip("\nUsage: remove <beer name>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   remove Sixpoint Resin"))
		fmt.Println(ui.Info("   remove \"Hazy IPA\""))
		return nil
	}

	// get the beer name from command arguments
	beerName := ui.ToTitleCase(strings.Join(cmd.Args, " "))

	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		fmt.Println(ui.Error("error loading favorites", err))
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
		fmt.Println(ui.Warning(" Beer '%s' is not in your favorites", beerName))
		fmt.Println(ui.Tip("\nPossible reasons:"))
		fmt.Println(ui.Info("   • The beer name might be spelled differently"))
		fmt.Println(ui.Info("   • It might have already been removed"))
		fmt.Println(ui.Tip("\nView your current favorites:"))
		fmt.Println(ui.Info("   favorites"))
		return nil
	}

	// remove the beer from the slice
	favorites.Beers = append(favorites.Beers[:index], favorites.Beers[index+1:]...)

	// save updated favorites
	err = storage.SaveFavorites(s.Config.User, favorites)
	if err != nil {
		fmt.Println(ui.Error("error saving favorites", err))
		return nil
	}
	fmt.Println(ui.Success("Beer %q removed from favorites!\n", beerName))
	return nil
}

// create a clear favorites command that removes all favorite beers
func clearFavoritesCommand(s *models.State, cmd models.Command) error {
	// create an empty favorites struct
	favorites := models.Favorites{Beers: []models.Beer{}}

	// ask are you sure if you want to clear favorites list
	fmt.Print(ui.Prompt("Are you sure you want to clear all your favorite beers? (y/n): "))
	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice != "y" && choice != "yes" {
		fmt.Println(ui.Info("Favorites not cleared."))
		return nil
	}

	// save the empty favorites to the file
	err := storage.SaveFavorites(s.Config.User, favorites)
	if err != nil {
		fmt.Println(ui.Error("error clearing favorites", err))
		return nil
	}

	fmt.Println(ui.Success("All favorite beers have been cleared."))
	return nil
}

// create a randomCommand that displays a random beer from last search results
func randomCommand(s *models.State, cmd models.Command) error {
	// check if ther are any last search results
	if len(s.LastSearchResults) == 0 {
		fmt.Println(ui.Warning(" No last search results found. Please perform a search first."))
		fmt.Println(ui.Tip("\nTo use the random feature:"))
		fmt.Println(ui.Info("   1. First search for beers: search IPA"))
		fmt.Println(ui.Info("   2. Then get a random suggestion: random"))
		fmt.Println(ui.Tip("\nTip: The random command picks from your last search"))
		return nil
	}

	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// generate a random index
	randomIndex := rand.Intn(len(s.LastSearchResults))

	// get the random beer
	randomBeer := s.LastSearchResults[randomIndex]

	// display the random beer information
	fmt.Println(ui.Header("\nRandom Beer from Last Search Results:\n\n"))

	// add Formatted beer display
	fmt.Println(ui.FormatBeer(randomBeer, 1))

	// ask user if they want to favorite the beer
	fmt.Print(ui.Prompt("Would you like to add this beer to your favorites? (y/n): "))
	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice == "y" || choice == "yes" {
		// create a command to favorite the beer
		favCmd := models.Command{
			Name: "favorite",
			Args: []string{randomBeer.Name},
		}
		// call the favorite command
		return favoriteCommand(s, favCmd)
	} else {
		fmt.Println(ui.Info("Beer not added to favorites."))
	}

	return nil
}

// create a history command function to display search history
func historyCommand(s *models.State, cmd models.Command) error {
	// check if there is any search history
	if len(s.SearchHistory) == 0 {
		fmt.Println(ui.Info("Your search history is empty"))
		fmt.Println(ui.Tip("\nSearch history will appear here after you:"))
		fmt.Println(ui.Info("   • Perform your first search"))
		fmt.Println(ui.Info("   • Example: search IPA"))
		fmt.Println(ui.Tip("\nTip: History is saved per user and persists between sessions"))
		return nil
	}

	// display search history
	fmt.Println(ui.Header("\nYour Search History:"))

	// iterate over search history entries and print their details
	for i := range s.SearchHistory {
		fmt.Println(ui.FormatSearchHistory(s.SearchHistory, i+1))
	}

	return nil
}

// add a clear history command function to clear search history
func clearHistoryCommand(s *models.State, cmd models.Command) error {
	// CREATE an empty history slice
	history := []models.SearchHistoryEntry{}

	// check if there is any history to clear
	if len(s.SearchHistory) == 0 {
		fmt.Println(ui.Warning("No search history to clear."))
		fmt.Println(ui.Tip("\nYour search history is already empty."))
		fmt.Println(ui.Info("   Perform searches to build your history."))
		return nil
	}

	// ask are you sure
	fmt.Print(ui.Prompt("Are you sure you want to clear your search history? (y/n): "))
	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice != "y" && choice != "yes" {
		fmt.Println(ui.Info("Search history not cleared."))
		return nil
	}
	// save the empty history to the file
	err := storage.SaveSearchHistory(s.Config.User, history)
	if err != nil {
		fmt.Println(ui.Error("error clearing search history", err))
		return nil
	}

	// clear in-memory history
	s.SearchHistory = []models.SearchHistoryEntry{}

	fmt.Println(ui.Success("All search history has been cleared."))
	return nil
}

// ---------------- Print Helper Functions ------------------ //

// createa a export favorites command function
func exportFavoritesCommand(s *models.State, cmd models.Command) error {
	// validate format argument
	if len(cmd.Args) == 0 {
		fmt.Println(ui.Error("Error: Export format is required (json/csv/txt)"))
		fmt.Println(ui.Tip("\nUsage: export favs <format>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   export favs json"))
		fmt.Println(ui.Info("   export favs csv"))
		fmt.Println(ui.Info("   export favs txt"))
		return nil
	}

	// load favorites from storage
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		fmt.Println(ui.Error("error loading favorites", err))
		return nil
	}

	// get appropriate exporter using factory
	format := strings.ToLower(cmd.Args[0])
	var exporter export.Exporter

	switch format {
	case "json":
		exporter = export.NewJSONExporter()
	case "csv":
		exporter = export.NewCSVExporter()
	case "txt":
		exporter = export.NewTXTExporter()
	default:
		fmt.Println(ui.Error("Error: Unsupported export format %q. Use 'json', 'csv', or 'txt'.", format))
		return nil
	}

	// generate filename
	filename := export.GenerateFilename(s.Config.User, "favorites", format)

	// call the exporter.Export() method
	err = exporter.Export(favorites.Beers, filename)
	if err != nil {
		fmt.Println(ui.Error("Error: Failed to export favorites: %v", err))
		return nil
	}

	// display success msg with file location
	fmt.Println(ui.Success("Favorites exported successfully to %s", filename))
	return nil
}

func main() {

	// prompt the user to login or continue as guest

	username := ""
	fmt.Println(ui.Prompt("Enter your username (or press Enter to continue as guest): "))
	fmt.Scanln(&username)

	// welcome message
	ui.ShowWelcomeBanner(username)

	// print app header
	fmt.Println(ui.StarDivider())

	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(ui.Error("Error loading .env file"))
		return
	}

	// get API key from .env file
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		ui.ShowAPIKeyError()
		return
	}

	// create a config instance
	cfg := &models.Config{
		APIKey: apiKey,
		User:   username,
	}

	// load search history for the user
	history, err := storage.LoadSearchHistory(username)
	if err != nil {
		// if file not found, initialize empty history
		if os.IsNotExist(err) {
			history = []models.SearchHistoryEntry{}
		} else {
			fmt.Println(ui.Warning("Warning: Could not load search history: %v\n", err))
			history = []models.SearchHistoryEntry{}
		}
	}

	// create a models.State instance
	s := &models.State{
		Config:        cfg,
		SearchHistory: history,
	}

	// create a commandHandler instance
	ch := &models.CommandHandler{
		Handlers: make(map[string]func(*models.State, models.Command) error),
	}

	// register commands
	ch.RegisterCommand("logout", logoutCommand)
	ch.RegisterCommand("login", loginCommand)
	ch.RegisterCommand("search", searchCommand)
	ch.RegisterCommand("help", helpCommand)
	ch.RegisterCommand("exit", exitCommand)
	ch.RegisterCommand("quit", exitCommand)
	ch.RegisterCommand("favorite", favoriteCommand)
	ch.RegisterCommand("favorites", displayFavoritesCommand)
	ch.RegisterCommand("remove", removeFavoriteCommand)
	ch.RegisterCommand("clear favs", clearFavoritesCommand)
	ch.RegisterCommand("random", randomCommand)
	ch.RegisterCommand("history", historyCommand)
	ch.RegisterCommand("clear history", clearHistoryCommand)
	ch.RegisterCommand("export favs", exportFavoritesCommand)

	// CLI interaction section

	// Main Loop
	for {

		// display current user
		if s.Config.User != "" {
			fmt.Println(ui.Header("Current User: %s\n", ui.ToTitleCase(s.Config.User)))
		} else {
			fmt.Println(ui.Header("Current User: guest"))
		}

		// prompt user for command
		fmt.Println(ui.Prompt("Enter a command (type 'help' for available commands): "))
		reader := bufio.NewReader(os.Stdin)
		// read user input
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// split input into command name and arguments
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue // skip empty input
		}

		// parse input into command struct
		// handle multi-word commands like "clear favs", "clear history", and "export favs"
		var cmd models.Command
		if len(parts) >= 2 && parts[0] == "clear" && (parts[1] == "favs" || parts[1] == "history") || len(parts) >= 2 && parts[0] == "export" && parts[1] == "favs" {

			// multi-word command
			cmd = models.Command{
				Name: parts[0] + " " + parts[1], // combine first two parts
				Args: parts[2:],                 // remaining parts are arguments
			}
		} else {
			// single-word command
			cmd = models.Command{
				Name: parts[0],  // first part is command name
				Args: parts[1:], // remaining parts are arguments
			}
		}

		// run the command
		err := ch.RunCommand(s, cmd)
		if err != nil {
			fmt.Println(ui.Error("Error executing command:", err))
			fmt.Println(ui.Tip("Type 'help' to see a full list of available commands."))
		}

		// end of main loop
		fmt.Println()

	}

}
