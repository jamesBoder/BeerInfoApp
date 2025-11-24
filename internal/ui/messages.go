package ui

import (
	"fmt"
	"os"
)

// ShowWelcomeBanner displays a welcome banner with an optional username
func ShowWelcomeBanner(username string) {
	banner := `
 __      __   _                    _____      _   _          ___                ___       __         _
 \ \    / /__| |__ ___ _ __  ___  |_   _|__  | |_| |_  ___  | _ ) ___ ___ _ _  |_ _|_ _  / _|___    /_\  _ __ _ __
  \ \/\/ / -_) / _/ _ \ '  \/ -_)  | |/ _ \ |  _| ' \/ -_) | _ \/ -_) -_) '_|  | || ' \|  _/ _ \  / _ \| '_ \ '_ \
   \_/\_/\___|_\__\___/_|_|_\___|  |_|\___/  \__|_||_\___| |___/\___\___|_|   |___|_||_|_| \___/ /_/ \_\ .__/ .__/
                                                                                                         |_|  |_|`

	fmt.Println(Header(banner))
	fmt.Println(StarDivider())
	fmt.Println()

	if username != "" {
		fmt.Println(Header(fmt.Sprintf("Hello there, %s!", ToTitleCase(username))))
	}

	// create an extra line space after hello message
	fmt.Println()

}

// display APIKEY ERROR message
func ShowAPIKeyError() {
	fmt.Println(Error("Error: API key not found in .env file"))
	fmt.Println(Warning("Checked file: .env"))
	fmt.Println()
	fmt.Println(Info("How to fix:"))
	fmt.Println(Info("   1. Open the .env file"))
	fmt.Println(Info("   2. Make sure it contains:"))
	fmt.Println(Info("      API_KEY=your_actual_key_here"))
	fmt.Println(Info("   3. Save the file and restart the app"))
	fmt.Println()
	fmt.Println(Warning("🔑 Get an API key:"))
	fmt.Println(Info("   https://rapidapi.com/winevybe/api/beer9"))
	fmt.Println()
	fmt.Println(Warning("The app will now exit"))
	ShowGoodbyeMessage()
}

// ShowGoodbyeMessage displays a goodbye message
func ShowGoodbyeMessage() {
	fmt.Println()
	fmt.Println(StarDivider())
	fmt.Println(Header("Thank you for using BeerInfoApp! 🍺"))
	fmt.Println(StarDivider())
	fmt.Println(Info(" Goodbye! Have a great day!"))
	os.Exit(0)
}
