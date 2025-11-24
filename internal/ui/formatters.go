package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// create FormatBeer() to format beer info for display
func FormatBeer(beer models.Beer, number int) string {
	// create header
	header := BeerNumber(number)

	// add each field with label. add Divider() between fields.
	// Name, Brewery, SKU,ABV,IBU,Category,Subcategory_1,Subcategory_2,Region,Country,Rating,Food Pairing,Description
	var sb strings.Builder
	sb.WriteString(header + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Name: " + beer.Name + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Brewery: " + beer.Brewery + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("SKU: " + beer.Sku + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("ABV: " + beer.Abv + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("IBU: " + beer.Ibu + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Category: " + beer.Category + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Subcategory 1: " + beer.SubCategory_1 + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Subcategory 2: " + beer.SubCategory_2 + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Region: " + beer.Region + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Country: " + beer.Country + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Rating: " + beer.Rating + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Food Pairing: " + beer.FoodPairing + "\n")
	sb.WriteString(Divider() + "\n")
	sb.WriteString("Description: " + beer.Description + "\n")

	return sb.String()

}

// create FormatBeerSummary() to create beer summary. Name, Brewery, SKU,ABV,IBU,Category,Subcategory_1,Subcategory_2,Region,Country,Rating,Food Pairing,Description
func FormatBeerSummary(beer models.Beer, number int) string {
	var sb strings.Builder
	sb.WriteString(BeerNumber(number))
	sb.WriteString(Divider() + "\n")

	// add fields conditionally. check if item is not empty
	if beer.Name != "" {
		sb.WriteString(CyanBold("Name: " + beer.Name + "\n"))
	}
	if beer.Brewery != "" {
		sb.WriteString(("Brewery: " + beer.Brewery + "\n"))
	}
	if beer.Sku != "" {
		sb.WriteString(CyanBold("SKU: " + beer.Sku + "\n"))
	}
	if beer.Abv != "" {
		sb.WriteString("ABV: " + beer.Abv + "\n")
	}
	if beer.Ibu != "" {
		sb.WriteString(CyanBold("IBU: " + beer.Ibu + "\n"))
	}
	if beer.Category != "" {
		sb.WriteString(("Category: " + beer.Category + "\n"))
	}
	if beer.SubCategory_1 != "" {
		sb.WriteString(CyanBold("Subcategory 1: " + beer.SubCategory_1 + "\n"))
	}
	if beer.SubCategory_2 != "" {
		sb.WriteString("Subcategory 2: " + beer.SubCategory_2 + "\n")
	}
	if beer.Region != "" {
		sb.WriteString(CyanBold("Region: " + beer.Region + "\n"))
	}
	if beer.Country != "" {
		sb.WriteString("Country: " + beer.Country + "\n")
	}
	if beer.Rating != "" {
		sb.WriteString(CyanBold("Rating: " + beer.Rating + "\n"))
	}
	if beer.FoodPairing != "" {
		sb.WriteString("Food Pairing: " + beer.FoodPairing + "\n")
	}
	if beer.Description != "" {
		sb.WriteString(CyanBold("Description: " + beer.Description + "\n"))
	}

	return sb.String()
}

// create a FormatSearchHistory() to format beer history info for display
func FormatSearchHistory(entry []models.SearchHistoryEntry, number int) string {
	var sb strings.Builder

	// validate index (number is 1-based in your callers)
	if number <= 0 || number > len(entry) {
		return ""
	}
	e := entry[number-1]

	// header with search emoji. no divider
	sb.WriteString(SearchNumber(number) + "\n")

	// term (cyan bold) no divider
	sb.WriteString(CyanBold("Search: " + e.Term + "\n"))

	// timestamp (magenta) no divider
	sb.WriteString("Timestamp: " + e.Timestamp.Format(time.RFC1123) + "\n")

	// Results count (cyan bold)
	sb.WriteString(CyanBold(fmt.Sprintf("Results Found: %d", len(e.Results))) + "\n")
	return sb.String()
}

// create a FormatHelpMenu to display help menu use commands and descriptions provided below
// color.Magenta("\nAvailable commands:")
// color.Cyan("---------------------------------------------------")
// coloredText := color.New(color.FgCyan).SprintFunc()
// fmt.Println(coloredText("  Command                   Description"))
// fmt.Println(coloredText("  -------------             -----------"))
// fmt.Println("  login <username>   - Log in with the specified username")
// fmt.Println("  logout             - Log out of the current session")
// fmt.Println("  search <beername>  - Search for a specific beer")
// fmt.Println("  favorite <beername> - Add a beer to your favorites")
// fmt.Println("  favorites          - Display your favorite beers")
// fmt.Println("  remove <beername>  - Remove a beer from your favorites")
// fmt.Println("  clear favs         - Clear all favorite beers")
// fmt.Println("  history            - Display your search history")
// fmt.Println("  clear history      - Clear your search history")
// fmt.Println("  export favs <format> - Export your favorites (json/csv/txt)")
// fmt.Println("  help               - Show this help message")
// fmt.Println("  exit, quit         - Exit the application")
// fmt.Println("  random             - Display a random beer from last search results")
func FormatHelpMenu() string {
	var sb strings.Builder
	sb.WriteString(Header("\nAvailable commands:") + "\n")
	sb.WriteString(Separator() + "\n")
	sb.WriteString("  Command                   Description\n")
	sb.WriteString("  -------------             -----------\n")
	sb.WriteString("  login <username>   - Log in with the specified username\n")
	sb.WriteString("  logout             - Log out of the current session\n")
	sb.WriteString("  search <beername>  - Search for a specific beer\n")
	sb.WriteString("  favorite <beername> - Add a beer to your favorites\n")
	sb.WriteString("  favorites          - Display your favorite beers\n")
	sb.WriteString("  remove <beername>  - Remove a beer from your favorites\n")
	sb.WriteString("  random             - Display a random beer from last search results\n")
	sb.WriteString("  clear favs         - Clear all favorite beers\n")
	sb.WriteString("  history            - Display your search history\n")
	sb.WriteString("  clear history      - Clear your search history\n")
	sb.WriteString("  export favs <format> - Export your favorites (json/csv/txt)\n")
	sb.WriteString("  help               - Show this help message\n")
	sb.WriteString("  exit, quit         - Exit the application\n")

	return sb.String()

}
