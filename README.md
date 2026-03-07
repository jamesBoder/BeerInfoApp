# BeerNfo

## Description

BeerNfo is a colorful CLI application that helps you discover and manage your favorite beers. Search for beers by name, save your favorites, view your search history, and export your collection—all from your terminal using the Beer9 API.

## Motivation

I am a beer connoisseur and wanted to create an application to lookup and find more information on my favorite beers around the world.

## Quick Start

**Prerequisites:** Go 1.21+ and a RapidAPI account

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd BeerInfoApp
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure API key**
   - Sign up at [RapidAPI](https://rapidapi.com/)
   - Subscribe to [Beer9 API](https://rapidapi.com/winevybe/api/beer9)
   - Create a `.env` file with your API key:
     ```env
     API_KEY=your_rapidapi_key_here
     ```

4. **Run the application**
   ```bash
   go run main.go
   ```
   You will be prompted to enter a username at startup (or press Enter to continue as guest).

## Usage

**Available Commands:**

| Command | Description |
|---------|-------------|
| `search <name>` | Search for beers by name or brewery |
| `random` | Get a random beer from the last search results |
| `favorite <name>` | Add a beer to favorites (by name from last search or name only) |
| `favorites` | View all favorite beers |
| `remove <name>` | Remove a beer from favorites |
| `clear favs` | Clear all favorite beers |
| `history` | View search history |
| `clear history` | Clear search history |
| `export favs <format>` | Export favorites to file (json/csv/txt) |
| `login <username>` | Switch to a different user |
| `logout` | Log out the current user |
| `help` | Show all commands |
| `exit` / `quit` | Exit the application |

**Example session:**
```
Enter your username (or press Enter to continue as guest):
> james

> search IPA
Found 5 beers...

> favorite Berkshire IPA
Beer "Berkshire Ipa" added to favorites with all details

> favorites
Your Favorite Beers:
1. Berkshire Ipa ...

> export favs json
Favorites exported successfully to james_favorites_<date>.json

> history
Your Search History:
1. "IPA" — 5 results

> exit
```

## Data Storage

User data is stored in JSON files in the working directory:
- **Favorites:** `favorites_<username>.json`
- **Search history:** `search_history_<username>.json`

Data persists between sessions and is scoped per user.

## Project Structure

```
BeerInfoApp/
├── main.go                    # Entry point: prompts for username, loads config, starts app
├── cmd/                       # Command implementations
│   ├── commands.go            # Registers all commands with the router
│   ├── auth.go                # login, logout commands
│   ├── search.go              # search, random commands
│   ├── favorites.go           # favorites, favorite, remove, clear favs commands
│   ├── history.go             # history, clear history commands
│   ├── export.go              # export favs command
│   └── help.go                # help, exit commands
├── internal/
│   ├── api/                   # Beer9 API client
│   ├── app/                   # App lifecycle, config loading, username prompt
│   ├── export/                # JSON, CSV, TXT exporters
│   ├── models/                # Shared data structs (Beer, State, Favorites, etc.)
│   ├── storage/               # Favorites and search history persistence
│   └── ui/                    # Colors, formatters, prompts, messages
└── pkg/cli/                   # CLI router, command parser
```

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the issues page or submit a pull request.

---

**Author:** James 'boder' Macean
