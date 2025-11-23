package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
)

// create a helper function that returns the favorities filename for a given user
func getUserFavoritesFilename(username string) string {
	// if username is empty, use "guest"
	if username == "" {
		username = "guest"
	}
	// trim spaces and convert to lowercase
	username = strings.TrimSpace(strings.ToLower(username))
	// replace spaces with underscores
	safeUsername := strings.ReplaceAll(username, " ", "_")
	// return the filename
	return fmt.Sprintf("favorites_%s.json", safeUsername)
}

// create a saveFavorites function to save favorite beers to a JSON file
func SaveFavorites(username string, favorites models.Favorites) error {
	// get the filename for the user
	filename := getUserFavoritesFilename(username)
	// create or truncate the favorites file
	file, err := os.Create(filename)
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
func LoadFavorites(username string) (models.Favorites, error) {
	var favorites models.Favorites

	// get username
	filename := getUserFavoritesFilename(username)

	// open the favorites file
	file, err := os.Open(filename)
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

// create a saveSearchHistory function to persist search history to a JSON file
func SaveSearchHistory(username string, history []models.SearchHistoryEntry) error {
	// get the filename for the user
	filename := fmt.Sprintf("search_history_%s.json", strings.ToLower(strings.ReplaceAll(username, " ", "_")))
	// create or truncate the search history file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// encode the search history slice to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(history)
	if err != nil {
		return err
	}

	return nil
}

// create a loadSearchHistory function to load search history from a JSON file
func LoadSearchHistory(username string) ([]models.SearchHistoryEntry, error) {
	var history []models.SearchHistoryEntry

	// get username
	filename := fmt.Sprintf("search_history_%s.json", strings.ToLower(strings.ReplaceAll(username, " ", "_")))

	// open the search history file
	file, err := os.Open(filename)
	if err != nil {
		return history, err
	}
	defer file.Close()

	// decode the JSON data into the history slice
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	if err != nil {
		return history, err
	}

	return history, nil
}
