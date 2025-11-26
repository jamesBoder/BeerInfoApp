package cmd

import (
	"fmt"
	"os"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// LoginCommand - implements cli.Command
type LoginCommand struct {
	ctx *CommandContext
}

// LogoutCommand - implements cli.Command
type LogoutCommand struct {
	ctx *CommandContext
}

// NewLoginCommand returns a instance of LoginCommand
func NewLoginCommand(ctx *CommandContext) *LoginCommand {
	return &LoginCommand{ctx: ctx}
}

// name method
func (c *LoginCommand) Name() string {
	return "login"
}

// Description method
func (c *LoginCommand) Description() string {
	return "Log in with a username"
}

// usage method
func (c *LoginCommand) Usage() string {
	return "login <username>"
}

// Execute method
func (c *LoginCommand) Execute(args []string) error {
	// check if username argument is provided
	if len(args) == 0 {
		fmt.Println(ui.Error("Error: Username is required"))
		fmt.Println(ui.Tip("\nUsage: login <username>"))
		fmt.Println(ui.Example("\nExamples:"))
		fmt.Println(ui.Info("   login james"))
		fmt.Println(ui.Info("   login sarah"))
		fmt.Println(ui.Tip("\nTip: Press Enter at startup to continue as guest"))
		return nil
	}

	// get the username from the command arguments
	username := args[0]
	// set the username in the config
	c.ctx.State.Config.User = username

	// load search history
	history, err := storage.LoadSearchHistory(username)
	if err != nil {
		// if file not found, initialize empty history
		if os.IsNotExist(err) {
			c.ctx.State.SearchHistory = []models.SearchHistoryEntry{}
		} else {
			// other errors
			fmt.Println(ui.Error("error loading search history", err))
			return nil
		}
	} else {
		c.ctx.State.SearchHistory = history
	}

	// print a success message
	fmt.Println(ui.Success("User %s logged in successfully\n", username))
	return nil
}

// NewLogoutCommand returns a instance of LogoutCommand
func NewLogoutCommand(ctx *CommandContext) *LogoutCommand {
	return &LogoutCommand{ctx: ctx}
}

// name method
func (c *LogoutCommand) Name() string {
	return "logout"
}

// Description method
func (c *LogoutCommand) Description() string {
	return "Log out the current user"
}

// usage method
func (c *LogoutCommand) Usage() string {
	return "logout"
}

// Execute method
func (c *LogoutCommand) Execute(args []string) error {
	// check if a user is logged in
	if c.ctx.State.Config.User == "" {
		fmt.Println(ui.Tip("No user is currently logged in"))
		return nil
	}

	// clear the username in the config
	c.ctx.State.Config.User = ""

	// clear search history from state
	c.ctx.State.SearchHistory = []models.SearchHistoryEntry{}

	// print a success message
	fmt.Println(ui.Success("User logged out successfully\n"))
	return nil
}
