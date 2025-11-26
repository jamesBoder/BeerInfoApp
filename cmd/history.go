package cmd

import (
	"fmt"

	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// HistoryCommand - implements cli.Command
type HistoryCommand struct {
	ctx *CommandContext
}

// ClearHistoryCommand - implements cli.Command
type ClearHistoryCommand struct {
	ctx *CommandContext
}

// NewHistoryCommand returns a instance of HistoryCommand
func NewHistoryCommand(ctx *CommandContext) *HistoryCommand {
	return &HistoryCommand{ctx: ctx}
}

// NewClearHistoryCommand returns a instance of ClearHistoryCommand
func NewClearHistoryCommand(ctx *CommandContext) *ClearHistoryCommand {
	return &ClearHistoryCommand{ctx: ctx}
}

// name, description, usage, and execute methods for each command will go here
// Implement the methods for HistoryCommand and ClearHistoryCommand as needed.

// HistoryCommand methods
func (c *HistoryCommand) Name() string {
	return "history"
}

func (c *HistoryCommand) Description() string {
	return "List search history"
}

func (c *HistoryCommand) Usage() string {
	return "history"
}

func (c *HistoryCommand) Execute(args []string) error {
	// check if there is any search history
	if len(c.ctx.State.SearchHistory) == 0 {
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
	for i := range c.ctx.State.SearchHistory {
		fmt.Println(ui.FormatSearchHistory(c.ctx.State.SearchHistory, i+1))
	}

	return nil
}

// ClearHistoryCommand methods
func (c *ClearHistoryCommand) Name() string {
	return "clear history"
}

func (c *ClearHistoryCommand) Description() string {
	return "Clear search history"
}

func (c *ClearHistoryCommand) Usage() string {
	return "clear history"
}

func (c *ClearHistoryCommand) Execute(args []string) error {
	// CREATE an empty history slice
	history := []models.SearchHistoryEntry{}

	// check if there is any history to clear
	if len(c.ctx.State.SearchHistory) == 0 {
		fmt.Println(ui.Warning("No search history to clear."))
		fmt.Println(ui.Tip("\nYour search history is already empty."))
		fmt.Println(ui.Info("   Perform searches to build your history."))
		return nil
	}

	// save the empty history to the file
	err := storage.SaveSearchHistory(c.ctx.State.Config.User, history)
	if err != nil {
		fmt.Println(ui.Error("error clearing search history", err))
		return nil
	}

	// clear in-memory history
	c.ctx.State.SearchHistory = []models.SearchHistoryEntry{}

	fmt.Println(ui.Success("All search history has been cleared."))
	return nil
}
