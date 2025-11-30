package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jamesBoder/BeerInfoApp.git/cmd"
	"github.com/jamesBoder/BeerInfoApp.git/internal/models"
	"github.com/jamesBoder/BeerInfoApp.git/internal/storage"
	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"
	"github.com/jamesBoder/BeerInfoApp.git/pkg/cli"
)

// Main application struct and lifecycle management

// create app struct
type App struct {
	config *Config
	state  *models.State
	router *cli.Router
	ctx    *cmd.CommandContext
}

// create a constructor for the app
func NewApp(config *Config) (*App, error) {
	// load search history for the user
	history, err := storage.LoadSearchHistory(config.Username)
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
	state := &models.State{
		Config:        &models.Config{APIKey: config.APIKey, User: config.Username},
		SearchHistory: history,
	}

	// create CLI router
	router := cli.NewRouter()

	// create command context
	ctx := &cmd.CommandContext{
		State: state,
	}

	// register all commands
	cmd.RegisterAllCommands(router, ctx)

	// register multi-world commands
	router.RegisterMultiWordCommand("clear favs")
	router.RegisterMultiWordCommand("clear history")
	router.RegisterMultiWordCommand("export favs")

	// create app instance
	app := &App{
		config: config,
		state:  state,
		router: router,
		ctx:    ctx,
	}

	return app, nil
}

// run() method to start the app. Shows welcom banner and starts the main loop, handles user input and executes commands, return error if fatal
func (app *App) Run() error {
	// show welcome banner
	ui.ShowWelcomeBanner(app.config.Username)
	fmt.Println(ui.StarDivider())

	// start main loop
	for {
		// display current user
		if app.state.Config.User != "" {
			fmt.Printf(ui.Info("Current User: %s\n", ui.ToTitleCase(app.state.Config.User)))
		} else {
			fmt.Println(ui.Info("Current User: guest"))
		}

		// prompt user for command
		fmt.Print(ui.Prompt("Enter a command (type 'help' for available commands): "))

		// read user input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// skip empty input
		if input == "" {
			continue
		}

		// parse and execute command
		err := app.router.Execute(input)
		if err != nil {
			if cmdErr, ok := err.(*cli.CommandError); ok {
				fmt.Println(ui.Error("Error: %s", cmdErr.Command))
				fmt.Println(ui.Tip("Type 'help' for available commands"))
			} else {
				fmt.Println(ui.Error("An error occurred: %v", err))
			}
		}
	}

	return nil
}

// Shutdown method that saves any pending state and cleans up resources, returns error if any issues occur
func (app *App) Shutdown() error {
	// save search history
	err := storage.SaveSearchHistory(app.config.Username, app.state.SearchHistory)
	if err != nil {
		return fmt.Errorf("failed to save search history: %v", err)
	}

	return nil
}
