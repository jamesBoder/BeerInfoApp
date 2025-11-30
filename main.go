package main

import (
	"fmt"

	"github.com/jamesBoder/BeerInfoApp.git/internal/ui"

	"github.com/jamesBoder/BeerInfoApp.git/internal/app"
)

// build command functions here

func main() {

	// prompt for user name using app.PromptForUsername()
	username := app.PromptForUsername()

	// load config using app.LoadConfig()
	config, err := app.LoadConfig()
	if err != nil {
		fmt.Println(ui.Error("Error loading config: %v\n", err))
		return
	}

	// set username in config
	config.Username = username

	// create new app instance using NewApp()
	application, err := app.NewApp(config)
	if err != nil {
		fmt.Println(ui.Error("Error initializing app: %v\n", err))
		return
	}

	// call application.Run()()
	application.Run()

	// CLI interaction section

	// Main Loop
	for {

		// end of main loop
		fmt.Println()

	}

}
