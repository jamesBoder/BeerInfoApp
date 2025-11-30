# 🍺 Beer Info App

A colorful CLI application for discovering and managing your favorite beers using the Beer9 API.

## ✨ Features

- 🔍 **Search Beers** - Search for beers by name with detailed information
- 🎲 **Random Beer** - Get random beer suggestions from your search results
- ⭐ **Favorites System** - Save and manage your favorite beers with complete details
- 👤 **Multi-User Support** - Each user has their own favorites and data
- 📊 **Search History** - Track all your searches with timestamps
- 📤 **Export Favorites** - Export your favorites to JSON, CSV, or TXT files
- 🎨 **Colorful Interface** - Beautiful color-coded output for better readability
- 💾 **Data Persistence** - Your favorites and history are saved between sessions

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher
- RapidAPI account with Beer9 API access

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd BeerInfoApp
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the project root:
```env
API_KEY=your_rapidapi_key_here
```

4. Run the application:
```bash
go run main.go
```

## 📖 Usage

### Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `search <beer name>` | Search for beers by name | `search IPA` |
| `random` | Get a random beer suggestion | `random` |
| `favorite <beer name>` | Add a beer to your favorites | `favorite Berkshire IPA` |
| `favorites` | Display all your favorite beers | `favorites` |
| `remove <beer name>` | Remove a beer from favorites | `remove Berkshire IPA` |
| `clear favs` | Clear all favorites | `clear favs` |
| `history` | View your search history | `history` |
| `clear history` | Clear your search history | `clear history` |
| `export favs <format>` | Export favorites (json/csv/txt) | `export favs json` |
| `login <username>` | Login as a user | `login james` |
| `logout` | Logout current user | `logout` |
| `help` | Show help message | `help` |
| `exit` | Exit the application | `exit` |

### Example Session

```bash
🍺 Welcome to the Beer Info App!
================================
Enter your username (or press Enter to continue as guest): james
Hello, james!

> search IPA
🍺 Found 5 beers:

--- Beer #1 ---
🍺 Name: Berkshire IPA
Brewery: Berkshire Brewing Company
ABV: 6.5%
...

> favorite Berkshire IPA
✅ Beer "Berkshire IPA" added to favorites!

> random
🎲 Random Beer Suggestion:
🍺 Name: Another Great IPA
...

> favorites
🍺 Your Favorite Beers (1 total):

--- Beer #1 ---
🍺 Name: Berkshire IPA
...

> export favs json
✅ Favorites exported successfully to favorites_james_20240115_143022.json

> logout
✅ User james logged out successfully
```

## 🔑 Getting Your API Key

1. Sign up at [RapidAPI](https://rapidapi.com/)
2. Subscribe to the [Beer9 API](https://rapidapi.com/winevybe/api/beer9)
3. Copy your API key from the dashboard
4. Add it to your `.env` file

## 📁 Project Structure

```
BeerInfoApp/
├── main.go                 # Entry point
├── cmd/                    # Command handlers
├── internal/               # Internal packages (models, storage, API, UI)
├── pkg/                    # Reusable CLI framework
├── .env                    # API key configuration
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
├── favorites_*.json        # User favorites (auto-generated)
├── search_history_*.json   # Search history (auto-generated)
└── README.md               # This file
```

## 🎨 Color Scheme

- 🟢 **Green** - Success messages
- 🔴 **Red** - Error messages
- 🟡 **Yellow** - Warnings
- 🔵 **Cyan** - Headers and beer names
- ⚪ **White** - Regular text and prompts

## 💾 Data Storage

Each user's data is stored in separate JSON files:

**Favorites:**
- `favorites_james.json` - James's favorites
- `favorites_sarah.json` - Sarah's favorites
- `favorites_guest.json` - Guest user favorites

**Search History:**
- `search_history_james.json` - James's search history
- `search_history_sarah.json` - Sarah's search history
- `search_history_guest.json` - Guest search history

Your data persists between sessions and survives logout/login.

## 🛠️ Built With

- [Go](https://golang.org/) - Programming language
- [Beer9 API](https://rapidapi.com/winevybe/api/beer9) - Beer data
- [fatih/color](https://github.com/fatih/color) - Terminal colors
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variables

## 📝 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

## 👤 Author

James 'boder' Macean

## 🙏 Acknowledgments

- Beer9 API for providing beer data
- RapidAPI for API hosting
- The Go community for excellent packages
- You for trying out this application
