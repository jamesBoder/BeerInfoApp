package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// favorite command implementation will go here

// FavoritesCommand - implements cli.Command
type FavoritesCommand struct {
	ctx *CommandContext
}

// FavoriteCommand - implements cli.Command
type FavoriteCommand struct {
	ctx *CommandContext
}

// RemoveCommand - implements cli.Command
type RemoveCommand struct {
	ctx *CommandContext
}

// ClearFavsCommand - implements cli.Command
type ClearFavsCommand struct {
	ctx *CommandContext
}

// NewFavoritesCommand returns a instance of FavoritesCommand
func NewFavoritesCommand(ctx *CommandContext) *FavoritesCommand {
	return &FavoritesCommand{ctx: ctx}
}

// NewFavoriteCommand returns a instance of FavoriteCommand
func NewFavoriteCommand(ctx *CommandContext) *FavoriteCommand {
	return &FavoriteCommand{ctx: ctx}
}

// NewRemoveCommand returns a instance of RemoveCommand
func NewRemoveCommand(ctx *CommandContext) *RemoveCommand {
	return &RemoveCommand{ctx: ctx}
}

// NewClearFavsCommand returns a instance of ClearFavsCommand
func NewClearFavsCommand(ctx *CommandContext) *ClearFavsCommand {
	return &ClearFavsCommand{ctx: ctx}
}

// name, description, usage, and execute methods for each command will go here
// Implement the methods for FavoritesCommand, FavoriteCommand, RemoveCommand, and ClearFavsCommand as needed.

// FavoritesCommand methods
func (c *FavoritesCommand) Name() string {
	return "favorites"
}

func (c *FavoritesCommand) Description() string {
	return "List all favorite beers"
}

func (c *FavoritesCommand) Usage() string {
	return "favorites"
}

func (c *FavoritesCommand) Execute(args []string) error {
	// load existing favorites
	favorites, err := storage.LoadFavorites(c.ctx.State.Config.User)
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

// FavoriteCommand methods
func (c *FavoriteCommand) Name() string {
	return "favorite"
}

func (c *FavoriteCommand) Description() string {
	return "Add a beer to favorites"
}

func (c *FavoriteCommand) Usage() string {
	return "favorite <beer id>"
}

func (c *FavoriteCommand) Execute(args []string) error {
	// check if beer name is provided
	if len(args) == 0 {
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
	beerName := ui.ToTitleCase(strings.Join(args, " "))

	// load existing favorites
	favorites, err := storage.LoadFavorites(c.ctx.State.Config.User)
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

	if len(c.ctx.State.LastSearchResults) > 0 {
		for _, beer := range c.ctx.State.LastSearchResults {
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
	err = storage.SaveFavorites(c.ctx.State.Config.User, favorites)
	if err != nil {
		fmt.Println(ui.Error("Error saving favorites:", err))
		return nil
	}

	return nil
}

// Remove Command methods
func (c *RemoveCommand) Name() string {
	return "remove"
}

func (c *RemoveCommand) Description() string {
	return "Remove a beer from favorites"
}

func (c *RemoveCommand) Usage() string {
	return "remove <beer id>"
}

func (c *RemoveCommand) Execute(args []string) error {
	// check if beer name is provided
	if len(args) == 0 {
		fmt.Println(ui.Warning("beer name not provided"))
		fmt.Println(ui.Tip("\nUsage: remove <beer name>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   remove Sixpoint Resin"))
		fmt.Println(ui.Info("   remove \"Hazy IPA\""))
		return nil
	}

	// get the beer name from command arguments
	beerName := ui.ToTitleCase(strings.Join(args, " "))

	// load existing favorites
	favorites, err := storage.LoadFavorites(c.ctx.State.Config.User)
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
	err = storage.SaveFavorites(c.ctx.State.Config.User, favorites)
	if err != nil {
		fmt.Println(ui.Error("error saving favorites", err))
		return nil
	}
	fmt.Println(ui.Success("Beer %q removed from favorites!\n", beerName))
	return nil
}

// ClearFavs Command methods
func (c *ClearFavsCommand) Name() string {
	return "clear favs"
}

func (c *ClearFavsCommand) Description() string {
	return "Clear all favorite beers"
}

func (c *ClearFavsCommand) Usage() string {
	return "clear favs"
}

func (c *ClearFavsCommand) Execute(args []string) error {
	// create an empty favorites struct
	favorites := models.Favorites{Beers: []models.Beer{}}

	// save the empty favorites to the file
	err := storage.SaveFavorites(c.ctx.State.Config.User, favorites)
	if err != nil {
		fmt.Println(ui.Error("error clearing favorites", err))
		return nil
	}

	fmt.Println(ui.Success("All favorite beers have been cleared."))
	return nil
}
