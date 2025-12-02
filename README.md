# 🍺 BeerNfo

## Description

BeerNfo is a colorful CLI application that helps you discover and manage your favorite beers. Search for beers by name, save your favorites, and export your collection—all from your terminal using the Beer9 API.

## Motivation

I am a beer connoisseur and wanted to create an application to lookup and find more information on my favorite beers around the world.

## Quick Start

**Prerequisites:** Go 1.16+ and a RapidAPI account

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

## Usage

**Essential Commands:**

| Command | Description |
|---------|-------------|
| `search <name>` | Search for beers |
| `favorite <name>` | Add beer to favorites |
| `favorites` | View all favorites |
| `remove <name>` | Remove from favorites |
| `export favs <format>` | Export favorites (json/csv/txt) |
| `random` | Get random beer from search results |
| `login <username>` | Login as user |
| `help` | Show all commands |
| `exit` | Exit application |

**Example:**
```bash
> search IPA
🍺 Found 5 beers...

> favorite Berkshire IPA
✅ Beer "Berkshire IPA" added to favorites!

> favorites
🍺 Your Favorite Beers (1 total)...
```

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the issues page or submit a pull request.

---

**Author:** James 'boder' Macean
