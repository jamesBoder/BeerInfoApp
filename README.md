# BeerNfo

**A colorful CLI application for discovering, saving, and exporting your favorite beers — powered by the Beer9 API.**

[![CI](https://github.com/jamesBoder/BeerInfoApp/actions/workflows/ci.yml/badge.svg)](https://github.com/jamesBoder/BeerInfoApp/actions/workflows/ci.yml)
[![Release](https://github.com/jamesBoder/BeerInfoApp/actions/workflows/release.yml/badge.svg)](https://github.com/jamesBoder/BeerInfoApp/releases)

---

## Motivation

I am a beer connoisseur and I wanted to create a fun application to look up and find more information about beers. I also wanted to practice building a more complex Go application with a clean architecture, reusable components, and a nice terminal UI.

---

## Features

- **Beer Search** — Query the Beer9 API by name with colored, formatted results
- **Favorites Management** — Add, view, and remove favorites persisted to disk
- **Multi-Format Export** — Export your collection as JSON, CSV, or TXT
- **Search History** — Tracks every search, filterable and clearable
- **Multi-User Support** — Per-user favorites and history via username sessions
- **Colorful Terminal UI** — Styled output using `fatih/color`

---

## Getting Started

**Option 1 — Download a pre-built binary** (no Go required)

Grab the latest release for your OS from the [Releases page](https://github.com/jamesBoder/BeerInfoApp/releases), then:

```bash
./beerinfo
```

**Option 2 — Install with Go**

```bash
go install github.com/jamesBoder/BeerInfoApp@latest
```

**Option 3 — Run from source**

```bash
git clone https://github.com/jamesBoder/BeerInfoApp.git
cd BeerInfoApp
go mod download
go run main.go
```

**API Key (required for all options)**

Create a `.env` file in the same directory as the binary:
```env
API_KEY=your_rapidapi_key_here
```

> Get a free key at [RapidAPI → Beer9 API](https://rapidapi.com/winevybe/api/beer9)

---

## Commands

| Command | Description |
|---|---|
| `search <name>` | Search beers by name via the Beer9 API |
| `random` | Get a random beer from the last search results |
| `favorite <name>` | Add a beer to your favorites |
| `favorites` | List all saved favorites |
| `remove <name>` | Remove a beer from favorites |
| `clear favs` | Clear all favorites |
| `export favs <format>` | Export favorites: `json`, `csv`, or `txt` |
| `history` | View your search history |
| `clear history` | Clear all search history |
| `login <username>` | Switch to a different user profile |
| `help` | Show all available commands |
| `exit` | Exit the application |

**Example session:**
```
> login james
✅ Logged in as James

> search IPA
🍺 Found 5 beers...

> favorite Berkshire IPA
✅ Beer "Berkshire IPA" added to favorites!

> export favs csv
✅ Exported 1 beer(s) to james_favs_20240315_143022.csv
```

---

## Project Structure

```
BeerInfoApp/
├── main.go
├── cmd/                        # Command handler layer
│   ├── commands.go             # Command registration
│   ├── auth.go                 # login/logout
│   ├── search.go               # search, random
│   ├── favorites.go            # favorite, favorites, remove, clear favs
│   ├── history.go              # history, clear history
│   ├── export.go               # export favs
│   └── help.go                 # help, exit
├── internal/
│   ├── app/                    # App lifecycle: startup, run loop, shutdown
│   ├── models/                 # Shared data structs
│   ├── api/                    # Beer9 API HTTP client
│   ├── storage/                # JSON file persistence
│   ├── export/                 # Strategy pattern exporters (JSON/CSV/TXT)
│   └── ui/                     # Colors, prompts, formatters
└── pkg/cli/                    # Reusable CLI framework
    ├── command.go              # Command interface + BaseCommand
    └── router.go               # Input parsing + command dispatch
```

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Terminal Colors | [fatih/color](https://github.com/fatih/color) |
| Env Loading | [joho/godotenv](https://github.com/joho/godotenv) |
| Beer Data | [Beer9 via RapidAPI](https://rapidapi.com/winevybe/api/beer9) |
| Persistence | JSON file storage (per-user) |

---

## Testing

Unit tests cover the CLI framework, all three export formats, and the storage layer.

```bash
go test ./pkg/cli/... ./internal/export/... ./internal/storage/...
```

Tests run automatically on every push and pull request via GitHub Actions.

---

## Contributing

Contributions, issues, and feature requests are welcome. Feel free to open an issue or submit a pull request.

---

**Author:** James 'boder' Macean
