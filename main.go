package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"

	"os"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/jamesBoder/BeerInfoApp.git/internal/api"
	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
)

// build command functions here

// create a login command functino
func loginCommand(s *models.State, cmd models.Command) error {
	// check if username argument is provided
	if len(cmd.Args) == 0 {
		color.Red("❌ Error: Username is required")
		color.Yellow("\n💡 Usage: login <username>")
		color.Yellow("\n📖 Examples:")
		color.Cyan("   login james")
		color.Cyan("   login sarah")
		color.Yellow("\n💭 Tip: Press Enter at startup to continue as guest")
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
			color.Red("error loading search history", err)
			return nil
		}
	} else {
		s.SearchHistory = history
	}

	// print a success message
	color.Green("User %s logged in successfully\n", username)
	return nil
}

// create a logout command function
func logoutCommand(s *models.State, cmd models.Command) error {
	// check if user is logged in
	if s.Config.User == "" {
		color.Yellow("⚠️ No user is currently logged in.")
		color.Yellow("\n 💡 To log in, use the command: login <username>")
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
		color.Red("error saving guest favorites", err)
		return nil
	}

	// save empty search history for guest user
	err = storage.SaveSearchHistory("guest", []models.SearchHistoryEntry{})
	if err != nil {
		color.Red("error saving guest search history", err)
		return nil
	}

	// print a success message
	color.Green("User %s logged out successfully\n", username)
	color.Green("Your favorite beers are saved and will be available when you log back in.")
	return nil
}

// create a search command function
func searchCommand(s *models.State, cmd models.Command) error {
	// base case : check if search term is provided
	if len(cmd.Args) == 0 {
		color.Red("❌ Error: Search term is required")
		color.Yellow("\n💡 Usage: search <beer name or brewery>")
		color.Yellow("\n📖 Examples:")
		color.Cyan("   search IPA")
		color.Cyan("   search Sixpoint")
		color.Cyan("   search \"Hazy IPA\"")
		color.Yellow("\n💭 Tip: Use quotes for multi-word searches")
		return nil
	}

	// get the search term from command arguments
	var searchTerm string = cmd.Args[0]

	// get the API key from config
	//apiKey := ""
	//if s != nil && s.Config != nil {
	//	apiKey = s.Config.APIKey
	//}

	// exit if user types "quit" or "exit"
	if searchTerm == "quit" || searchTerm == "exit" {
		color.Green("Thanks for using Beer Info App! Goodbye!")
		return nil
	}

	// check if beer name is empty
	if searchTerm == "" {
		fmt.Println("Please enter a valid beer name.")
		return nil
	}

	// check if API key is set
	//if apiKey == "" {
	//color.Red("❌ Error: API key is not configured")
	//color.Yellow("\n🔍 The app needs an API key to search for beers")
	//color.Cyan("\n💡 How to fix:")
	//color.Cyan("   1. Create a file named '.env' in the app directory")
	//color.Cyan("   2. Add this line: API_KEY=your_api_key_here")
	//color.Cyan("   3. Replace 'your_api_key_here' with your actual RapidAPI key")
	//return nil
	//}

	// API interaction section

	// create a BeerAPI client
	client := api.NewBeerAPIClient(s.Config.APIKey)

	// call the SearchBeers method
	apiResponse, err := client.SearchBeers(searchTerm)
	if err != nil {
		color.Red("❌ Error: Could not connect to beer database")
		color.Yellow("🔍 Reason: %v", err)
		color.Cyan("\n💡 Possible solutions:")
		color.Cyan("   • Check your internet connection")
		color.Cyan("   • Verify your API key is still valid")
		color.Cyan("   • The API service might be temporarily down")
		color.Cyan("   • Try again in a few moments")
		color.Yellow("\n🌐 API Status: https://rapidapi.com/status")
		return nil
	}

	// Check if the API returned an error
	if apiResponse.Error {
		color.Red("❌ Error: Beer database returned an error")
		color.Yellow("🔢 Error code: %d", apiResponse.Code)

		// switch on error code

		switch apiResponse.Code {
		case 400:
			color.Red("❌ Error: Bad request")
			color.Yellow("🔍 Reason: Invalid search parameters")
		case 401:
			color.Red("❌ Error: Unauthorized")
			color.Yellow("🔍 Reason: Invalid API key or access denied")
		case 403:
			color.Red("❌ Error: Forbidden")
			color.Yellow("🔍 Reason: Access to the resource is forbidden")
		case 404:
			color.Red("❌ Error: Not found")
			color.Yellow("🔍 Reason: The requested resource was not found")
		case 500:
			color.Red("❌ Error: Internal server error")
			color.Yellow("🔍 Reason: The server encountered an error processing your request")
			color.Cyan("   • This is not your fault")
			color.Cyan("   • Try again in a few minutes")
		default:
			color.Red("❌ Error: Unknown error occurred (code %d)", apiResponse.Code)
		}

		return nil
	}

	if len(apiResponse.Data) == 0 {
		color.Yellow("🔍 No beers found matching '%s'", searchTerm)
		color.Cyan("\n💡 Suggestions:")
		color.Cyan("   • Try a broader search term")
		color.Cyan("   • Check your spelling")
		color.Cyan("   • Try searching by brewery name")
		color.Yellow("\n📖 Examples:")
		color.Cyan("   search IPA")
		color.Cyan("   search Sixpoint")
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
		color.Yellow("Warning: Could not save search history: %v\n", err)
	}

	// Print the decoded beer information
	color.Green("\n🍺 Found %d beers:\n\n", len(apiResponse.Data))
	for i, beer := range apiResponse.Data {
		// create bold yellow
		y := color.New(color.FgYellow, color.Bold)
		y.Printf("--- Beer #%d ---\n", i+1)
		fmt.Printf("Name: %s\n", beer.Name)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Brewery: %s\n", beer.Brewery)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("SKU: %s\n", beer.Sku)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("ABV: %s\n", beer.Abv)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("IBU: %s\n", beer.Ibu)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Category: %s\n", beer.Category)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Subcategory_1: %s\n", beer.SubCategory_1)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Subcategory_2: %s\n", beer.SubCategory_2)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Region: %s\n", beer.Region)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Country: %s\n", beer.Country)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Rating: %s\n", beer.Rating)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Food Pairing: %s\n", beer.FoodPairing)
		color.Cyan("---------------------------------------------------")
		fmt.Printf("Description: %s\n", beer.Description)
		color.Cyan("---------------------------------------------------")
		fmt.Println()
	}

	// Ask user what to do next

	color.Cyan("\nOptions:")
	color.Yellow("  [s] Search again")
	color.Blue("  [m] Main menu")
	color.Red("  [x] Exit")
	color.White("Your choice: ")

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
		color.Green("Thanks for using Beer Info App! Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("Returning to main menu...")
		return nil
	}
	return nil
}

// create a help command function
func helpCommand(s *models.State, cmd models.Command) error {
	color.Magenta("\nAvailable commands:")
	color.Cyan("---------------------------------------------------")
	coloredText := color.New(color.FgCyan).SprintFunc()
	fmt.Println(coloredText("  Command                   Description"))
	fmt.Println(coloredText("  -------------             -----------"))
	fmt.Println("  login <username>   - Log in with the specified username")
	fmt.Println("  logout             - Log out of the current session")
	fmt.Println("  search <beername>  - Search for a specific beer")
	fmt.Println("  random             - Display a random beer from last search results")
	fmt.Println("  favorite <beername> - Add a beer to your favorites")
	fmt.Println("  favorites          - Display your favorite beers")
	fmt.Println("  remove <beername>  - Remove a beer from your favorites")
	fmt.Println("  clear favs         - Clear all favorite beers")
	fmt.Println("  history            - Display your search history")
	fmt.Println("  clear history      - Clear your search history")
	fmt.Println("  export favs <format> - Export your favorites (json/csv)")
	fmt.Println("  help               - Show this help message")
	fmt.Println("  exit, quit         - Exit the application")
	return nil
}

// create an exit command function
func exitCommand(s *models.State, cmd models.Command) error {
	color.Green("Exiting the Beer Info App. Goodbye!")
	os.Exit(0)
	return nil
}

// create a title case function for beer names
func toTitleCase(input string) string {
	caser := cases.Title(language.English)
	return caser.String(input)
}

// create a favorite command function
func favoriteCommand(s *models.State, cmd models.Command) error {
	// check if beer name is provided
	if len(cmd.Args) == 0 {
		color.Red("❌ Error: Beer name is required")
		color.Yellow("\n💡 Usage: favorite <beer name>")
		color.Yellow("\n📖 Examples:")
		color.Cyan("   favorite Sixpoint Resin")
		color.Cyan("   favorite \"Hazy IPA\"")
		color.Yellow("\n💭 Tip: The beer should be from your last search results")
		color.Cyan("   Or type 'search <name>' first to find a beer")
		return nil
	}

	// get the beer name from command arguments
	beerName := toTitleCase(strings.Join(cmd.Args, " "))

	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		// if file not found, initialize empty favorites
		if os.IsNotExist(err) {
			favorites = models.Favorites{Beers: []models.Beer{}}
		} else {

			// other errors
			color.Red("Error loading favorites:", err)
			return nil
		}
	}

	// check if beer is already in favorites
	for _, beer := range favorites.Beers {
		if strings.EqualFold(beer.Name, beerName) {
			color.Yellow("Beer %q is already in your favorites.\n", beerName)
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
		color.Green("Beer %q added to favorites with all details\n", beerName)
	} else {
		// save a beer with only the name if not found in last search results
		// create empty beer with only name
		newBeer := models.Beer{Name: beerName}
		// append to favorites
		favorites.Beers = append(favorites.Beers, newBeer)
		color.Green("Beer %q added to favorites with name only\n", beerName)
	}

	// save updated favorites
	err = storage.SaveFavorites(s.Config.User, favorites)
	if err != nil {
		color.Red("Error saving favorites:", err)
		return nil
	}

	color.Yellow("Beer %q added to favorites!\n", beerName)
	return nil
}

// display favorite beers command function
func displayFavoritesCommand(s *models.State, cmd models.Command) error {
	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		color.Yellow("favorites is empty. Type 'help' to add a favorite beer", err)
		return nil
	}

	// check if there are any favorites
	if err != nil || len(favorites.Beers) == 0 {
		color.Yellow("📭 You don't have any favorite beers yet")
		color.Cyan("\n💡 How to add favorites:")
		color.Cyan("   1. Search for a beer: search IPA")
		color.Cyan("   2. Add to favorites: favorite <beer name>")
		color.Yellow("\n📖 Or type 'help' to see all commands")
		return nil
	}

	// display favorite beers
	color.Blue("\nYour Favorite Beers:")
	// iterate over favorite beers and print their names
	for i, beer := range favorites.Beers {
		d := color.New(color.FgCyan, color.Bold)
		color.Yellow("\n--- Beer #%d ---\n", i+1)
		d.Printf("Name: %s\n", beer.Name)

		// Only show fields if they exist
		if beer.Brewery != "" {
			fmt.Printf("Brewery: %s\n", beer.Brewery)
		}
		if beer.Abv != "" {
			d.Printf("ABV: %s\n", beer.Abv)
		}
		if beer.Ibu != "" {
			fmt.Printf("IBU: %s\n", beer.Ibu)
		}
		if beer.Category != "" {
			d.Printf("Category: %s\n", beer.Category)
		}
		if beer.SubCategory_1 != "" {
			fmt.Printf("Subcategory 1: %s\n", beer.SubCategory_1)
		}
		if beer.SubCategory_2 != "" {
			d.Printf("Subcategory 2: %s\n", beer.SubCategory_2)
		}
		if beer.Region != "" {
			fmt.Printf("Region: %s\n", beer.Region)
		}
		if beer.Country != "" {
			d.Printf("Country: %s\n", beer.Country)
		}
		if beer.Description != "" {
			fmt.Printf("Description: %s\n", beer.Description)
		}
		fmt.Println()
	}

	return nil
}

// add a remove favorite command function
func removeFavoriteCommand(s *models.State, cmd models.Command) error {
	// check if beer name is provided
	if len(cmd.Args) == 0 {
		color.Yellow("⚠️ beer name not provided")
		color.Cyan("\n💡 Usage: remove <beer name>")
		color.Yellow("\n📖 Examples:")
		color.Cyan("   remove Sixpoint Resin")
		color.Cyan("   remove \"Hazy IPA\"")
		return nil
	}

	// get the beer name from command arguments
	beerName := toTitleCase(strings.Join(cmd.Args, " "))

	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		color.Red("error loading favorites", err)
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
		color.Yellow("⚠️  Beer '%s' is not in your favorites", beerName)
		color.Cyan("\n💡 Possible reasons:")
		color.Cyan("   • The beer name might be spelled differently")
		color.Cyan("   • It might have already been removed")
		color.Yellow("\n📋 View your current favorites:")
		color.Cyan("   favorites")
		return nil
	}

	// remove the beer from the slice
	favorites.Beers = append(favorites.Beers[:index], favorites.Beers[index+1:]...)

	// save updated favorites
	err = storage.SaveFavorites(s.Config.User, favorites)
	if err != nil {
		color.Red("error saving favorites", err)
	}
	color.Green("Beer %q removed from favorites!\n", beerName)
	return nil
}

// create a clear favorites command that removes all favorite beers
func clearFavoritesCommand(s *models.State, cmd models.Command) error {
	// create an empty favorites struct
	favorites := models.Favorites{Beers: []models.Beer{}}

	// ask are you sure if you want to clear favorites list
	color.Cyan("Are you sure you want to clear all your favorite beers? (y/n): ")
	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice != "y" && choice != "yes" {
		color.Cyan("Favorites not cleared.")
		return nil
	}

	// save the empty favorites to the file
	err := storage.SaveFavorites(s.Config.User, favorites)
	if err != nil {
		color.Red("error clearing favorites", err)
		return nil
	}

	color.Green("All favorite beers have been cleared.")
	return nil
}

// create a randomCommand that displays a random beer from last search results
func randomCommand(s *models.State, cmd models.Command) error {
	// check if ther are any last search results
	if len(s.LastSearchResults) == 0 {
		color.Yellow("⚠️  No last search results found. Please perform a search first.")
		color.Cyan("\n💡 To use the random feature:")
		color.Cyan("   1. First search for beers: search IPA")
		color.Cyan("   2. Then get a random suggestion: random")
		color.Yellow("\n💭 Tip: The random command picks from your last search")
		return nil
	}

	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// generate a random index
	randomIndex := rand.Intn(len(s.LastSearchResults))

	// get the random beer
	randomBeer := s.LastSearchResults[randomIndex]

	// display the random beer information
	color.Green("\n🍺 Random Beer from Last Search Results:\n\n")
	y := color.New(color.FgYellow, color.Bold)
	y.Printf("--- Beer ---\n")
	fmt.Printf("Name: %s\n", randomBeer.Name)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Brewery: %s\n", randomBeer.Brewery)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("SKU: %s\n", randomBeer.Sku)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("ABV: %s\n", randomBeer.Abv)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("IBU: %s\n", randomBeer.Ibu)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Category: %s\n", randomBeer.Category)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Subcategory_1: %s\n", randomBeer.SubCategory_1)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Subcategory_2: %s\n", randomBeer.SubCategory_2)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Region: %s\n", randomBeer.Region)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Country: %s\n", randomBeer.Country)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Rating: %s\n", randomBeer.Rating)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Food Pairing: %s\n", randomBeer.FoodPairing)
	color.Cyan("---------------------------------------------------")
	fmt.Printf("Description: %s\n", randomBeer.Description)
	color.Cyan("---------------------------------------------------")
	fmt.Println()

	// ask user if they want to favorite the beer
	color.Cyan("Would you like to add this beer to your favorites? (y/n): ")
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
		color.Cyan("Beer not added to favorites.")
	}

	return nil
}

// create a history command function to display search history
func historyCommand(s *models.State, cmd models.Command) error {
	// check if there is any search history
	if len(s.SearchHistory) == 0 {
		color.Yellow("📭 Your search history is empty")
		color.Cyan("\n💡 Search history will appear here after you:")
		color.Cyan("   • Perform your first search")
		color.Cyan("   • Example: search IPA")
		color.Yellow("\n💭 Tip: History is saved per user and persists between sessions")
		return nil
	}

	// display search history
	color.Blue("\nYour Search History:")
	for i, entry := range s.SearchHistory {
		color.Yellow("\n--- Search #%d ---\n", i+1)
		fmt.Printf("Term: %s\n", entry.Term)
		fmt.Printf("Timestamp: %s\n", entry.Timestamp.Format(time.RFC1123))
		fmt.Printf("Results Found: %d\n", len(entry.Results))
	}
	return nil
}

// add a clear history command function to clear search history
func clearHistoryCommand(s *models.State, cmd models.Command) error {
	// CREATE an empty history slice
	history := []models.SearchHistoryEntry{}

	// check if there is any history to clear
	if len(s.SearchHistory) == 0 {
		color.Yellow("⚠️ No search history to clear.")
		color.Cyan("\n💡 Your search history is already empty.")
		color.Cyan("   Perform searches to build your history.")
		return nil
	}

	// ask are you sure
	color.Cyan("Are you sure you want to clear your search history? (y/n): ")
	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice != "y" && choice != "yes" {
		color.Cyan("Search history not cleared.")
		return nil
	}
	// save the empty history to the file
	err := storage.SaveSearchHistory(s.Config.User, history)
	if err != nil {
		color.Red("error clearing search history", err)
		return nil
	}

	// clear in-memory history
	s.SearchHistory = []models.SearchHistoryEntry{}

	color.Green("All search history has been cleared.")
	return nil
}

// ---------------- Print Helper Functions ------------------ //

// export beers to JSON function
func exportBeersToJSON(beers []models.Beer, filename string) error {
	// create or truncate the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// close the file when done
	defer file.Close()

	// create a JSON Encoder with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// encode the beers slice to JSON
	err = encoder.Encode(beers)
	if err != nil {
		return err
	}

	return nil
}

// export beers to CSV function
func exportBeersToCSV(beers []models.Beer, filename string) error {
	// create or truncate the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// close the file when done
	defer file.Close()

	// create a csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write the header row with all field names
	header := []string{"SKU", "Name", "Brewery", "Description", "Region", "Country", "ABV", "IBU", "Category", "Rating", "Food Pairing", "SubCategory_1", "SubCategory_2"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	// iterate over beers and write each as a row
	for _, beer := range beers {
		row := []string{
			beer.Sku,
			beer.Name,
			beer.Brewery,
			beer.Description,
			beer.Region,
			beer.Country,
			beer.Abv,
			beer.Ibu,
			beer.Category,
			beer.Rating,
			beer.FoodPairing,
			beer.SubCategory_1,
			beer.SubCategory_2,
		}
		err = writer.Write(row)
		if err != nil {
			return err
		}

	}
	return nil
}

// createa a export favorites command function
func exportFavoritesCommand(s *models.State, cmd models.Command) error {
	// check if format is provided
	if len(cmd.Args) == 0 {
		color.Red("❌ Error: Format is required")
		color.Yellow("\n💡 Usage: export favorites <format>")
		color.Yellow("\n📖 Examples:")
		color.Cyan("   export favorites json")
		color.Cyan("   export favorites csv")
		return nil
	}

	// get the format from the command arguments
	format := strings.ToLower(strings.TrimSpace(cmd.Args[0]))
	// check file extension

	if format != "json" && format != "csv" {
		color.Red("❌ Error: Unsupported format '%s'", format)
		color.Yellow("\n💡 Supported formats: json, csv")
		return nil
	}

	// load existing favorites
	favorites, err := storage.LoadFavorites(s.Config.User)
	if err != nil {
		color.Red("error loading favorites", err)
		color.Cyan("\n💡 Make sure you have favorite beers saved first.")
		return nil
	}

	// check if there are any favorites
	if len(favorites.Beers) == 0 {
		color.Yellow("📭 You don't have any favorite beers to export.")
		color.Cyan("\n💡 Add favorite beers first using the 'favorite' command.")
		return nil
	}

	// generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	safeUsername := strings.ToLower(strings.ReplaceAll(s.Config.User, " ", "_"))
	outputFilename := fmt.Sprintf("favorites_%s_%s.%s", safeUsername, timestamp, format)

	// call appropriate export function based on format
	if format == "json" {
		err = exportBeersToJSON(favorites.Beers, outputFilename)
	} else if format == "csv" {
		err = exportBeersToCSV(favorites.Beers, outputFilename)
	}
	if err != nil {
		color.Red("error exporting favorites", err)
		return nil
	}

	color.Green("✅ Successfully exported %d beers to %s", len(favorites.Beers), outputFilename)
	color.Yellow("\n📁 File location: %s", outputFilename)
	color.Cyan("\n💡 You can now:")
	color.Cyan("   • Open the file in a text editor")
	color.Cyan("   • Import it into a spreadsheet (for CSV)")
	color.Cyan("   • Share it with friends")

	return nil
}

func main() {

	color.Cyan(` __      __   _                    _____      _   _          ___                ___       __         _             
 \ \    / /__| |__ ___ _ __  ___  |_   _|__  | |_| |_  ___  | _ ) ___ ___ _ _  |_ _|_ _  / _|___    /_\  _ __ _ __ 
  \ \/\/ / -_) / _/ _ \ '  \/ -_)   | |/ _ \ |  _| ' \/ -_) | _ \/ -_) -_) '_|  | || ' \|  _/ _ \  / _ \| '_ \ '_ \
   \_/\_/\___|_\__\___/_|_|_\___|   |_|\___/  \__|_||_\___| |___/\___\___|_|   |___|_||_|_| \___/ /_/ \_\ .__/ .__/
                                                                                                        |_|  |_|   `)

	color.Cyan("--------------------------------------------------------------------------------------------------------------------")

	// prompt the user to login or continue as guest

	username := ""
	color.Magenta("Enter your username (or press Enter to continue as guest): ")
	fmt.Scanln(&username)

	// If username is empty, set to "Guest"
	if username == "" {
		username = "guest"
	}

	// greet the user
	color.Magenta("Hello there, %s!\n", toTitleCase(username))
	color.Cyan("********************************************************************************************************8*")

	// load .env file
	err := godotenv.Load()
	if err != nil {
		color.Red("Error loading .env file")
		return
	}

	// get API key from .env file
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		color.Red("❌ Error: API key not found in .env file")
		color.Yellow("🔍 Checked file: .env")
		color.Cyan("\n💡 How to fix:")
		color.Cyan("   1. Open the .env file")
		color.Cyan("   2. Make sure it contains:")
		color.Cyan("      API_KEY=your_actual_key_here")
		color.Cyan("   3. Save the file and restart the app")
		color.Yellow("\n🔑 Get an API key:")
		color.Cyan("   https://rapidapi.com/winevybe/api/beer9")
		color.Yellow("\n⚠️  The app will now exit")
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
			color.Yellow("Warning: Could not load search history: %v\n", err)
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
			color.Magenta("\nCurrent User: %s\n", toTitleCase(s.Config.User))
		} else {
			color.Magenta("\nCurrent User: guest")
		}

		// prompt user for command
		color.Magenta("\n> Enter a command (type 'help' for available commands): ")
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
		// handle multi-word commands like "clear favs" and "clear history"
		var cmd models.Command
		if len(parts) >= 2 && parts[0] == "clear" && (parts[1] == "favs" || parts[1] == "history") || len(parts) >= 2 && parts[0] == "export" && parts[1] == "favorites" {

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
			color.Red("Error executing command:", err)
			color.Blue("Type 'help' to see a full list of available commands.")
		}

		// end of main loop
		fmt.Println()

	}

}
