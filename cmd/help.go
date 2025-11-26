package cmd

import (
	"fmt"
	"os"

	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
)

// HelpCommand - implements cli.Command
type HelpCommand struct {
	ctx *CommandContext
}

// ExitCommand - implements cli.Command
type ExitCommand struct {
	ctx *CommandContext
}

// NewHelpCommand returns a instance of HelpCommand
func NewHelpCommand(ctx *CommandContext) *HelpCommand {
	return &HelpCommand{ctx: ctx}
}

// name method
func (c *HelpCommand) Name() string {
	return "help"
}

// Description method
func (c *HelpCommand) Description() string {
	return "Show help information"
}

// usage method
func (c *HelpCommand) Usage() string {
	return "help [command]"
}

// Execute method
func (c *HelpCommand) Execute(args []string) error {
	// use FormatHelpMenu to display help menu
	fmt.Println(ui.FormatHelpMenu())
	return nil
}

// NewExitCommand returns a instance of ExitCommand
func NewExitCommand(ctx *CommandContext) *ExitCommand {
	return &ExitCommand{ctx: ctx}
}

// name method
func (c *ExitCommand) Name() string {
	return "exit"
}

// Description method
func (c *ExitCommand) Description() string {
	return "Exit the application"
}

// usage method
func (c *ExitCommand) Usage() string {
	return "exit"
}

// Execute method
func (c *ExitCommand) Execute(args []string) error {
	ui.ShowGoodbyeMessage()
	os.Exit(0)
	return nil
}
