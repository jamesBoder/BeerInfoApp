package cmd

import (
	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/pkg/cli"
)

// CommandContext holds shared dependecies for all commands
type CommandContext struct {
	State *models.State
}

// RegisterAllCommands registers all comds with the routher
func RegisterAllCommands(router *cli.Router, ctx *CommandContext) {
	// Auth commands
	router.RegisterCommand(NewLoginCommand(ctx))
	router.RegisterCommand(NewLogoutCommand(ctx))

	// Search commands
	router.RegisterCommand(NewSearchCommand(ctx))
	router.RegisterCommand(NewRandomCommand(ctx))

	// Favorites commands
	router.RegisterCommand(NewFavoritesCommand(ctx))
	router.RegisterCommand(NewFavoriteCommand(ctx))
	router.RegisterCommand(NewClearFavsCommand(ctx))
	router.RegisterCommand(NewRemoveCommand(ctx))

	// History commands
	router.RegisterCommand(NewHistoryCommand(ctx))
	router.RegisterCommand(NewClearHistoryCommand(ctx))

	// Export commands
	router.RegisterCommand(NewExportFavsCommand(ctx))

	// Help command
	router.RegisterCommand(NewHelpCommand(ctx))
	router.RegisterCommand(NewExitCommand(ctx))

	// Register aliases
	router.RegisterAlias("quit", "exit")

	// Register mulit=word commands
	router.RegisterMultiWordCommand("clear favs")
	router.RegisterMultiWordCommand("clear history")
	router.RegisterMultiWordCommand("export favs")

}
