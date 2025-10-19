Note: I had this idea while out with friends, when I came home I just vibecoded the shit out of it to get something working before I lost the thought but I plan on coming back to this and actually doing it properly once I have a bit more time and adding a lot more functionality.

# TypeRacer TUI

A terminal-based TypeRacer clone built with Go using Charmbracelet's Bubble Tea, Lip Gloss, and Wish libraries.

## Features

- **Single-Player Practice Mode**: Practice typing with real-time WPM and accuracy tracking
- **Multiplayer Races**: SSH-based multiplayer typing races with real-time opponent progress
- **Beautiful TUI**: Styled terminal interface with color-coded typing feedback
- **Real-time Stats**: Live WPM calculation and accuracy tracking
- **Quote Integration**: Fetches random quotes from quotable.io API
- **Lobby System**: Matchmaking with configurable room sizes (2-4 players)
- **Countdown Timer**: 3-2-1-GO countdown before races start

## Installation

```bash
# Clone the repository
git clone https://github.com/sneak/typeracer-tui.git
cd typeracer-tui

# Build the application
go build -o typeracer-tui .
```

## Usage

### Practice Mode (Single Player)

```bash
# Run practice mode
./typeracer-tui
# or
./typeracer-tui -mode practice
```

### Server Mode (Multiplayer)

```bash
# Start the SSH server
./typeracer-tui -mode server

# With custom settings
./typeracer-tui -mode server -port 2222 -players 4
```

### Connecting to Server

```bash
# Connect via SSH
ssh localhost -p 2222

# Or connect to remote server
ssh user@server.com -p 2222
```

## Controls

- **Type**: Enter the displayed text as fast and accurately as possible
- **Backspace**: Correct mistakes
- **Ctrl+C / Esc**: Quit the application
- **r**: Restart (practice mode)
- **q**: Quit (results screen)

## Features

### Practice Mode
- Single-player typing practice
- Real-time WPM and accuracy tracking
- Visual feedback for correct/incorrect typing
- No network connection required

### Server Mode
- Multiplayer typing races over SSH
- Lobby system for player matchmaking
- Real-time opponent progress tracking
- Configurable room sizes (2-4 players)
- 3-2-1-GO countdown before races

### Visual Design
- Color-coded typing feedback (green for correct, red for errors)
- Progress bars and WPM displays
- Stylized countdown timer
- Race completion celebration
- Opponent progress visualization

## Architecture

The application is built using:

- **Bubble Tea**: TUI framework for managing application state and user input
- **Lip Gloss**: Terminal styling library for beautiful UI components
- **Wish**: SSH server framework for multiplayer functionality
- **Charmbracelet SSH**: SSH library for handling connections

### Project Structure

```
typeracer-tui/
├── main.go                 # Entry point with CLI flags
├── server.go              # SSH server setup with Wish
├── game/
│   ├── manager.go         # Game session & lobby management
│   ├── session.go         # Individual game session state
│   └── player.go          # Player state and progress
├── quotes/
│   └── fetcher.go         # Quote API integration
├── ui/
│   ├── practice.go        # Single-player Bubble Tea model
│   ├── multiplayer.go     # Multiplayer Bubble Tea model
│   ├── lobby.go           # Lobby waiting screen model
│   └── styles.go          # Lip Gloss styles
└── go.mod
```

## Development

### Prerequisites

- Go 1.24 or later
- SSH client for testing multiplayer mode

### Building

```bash
# Install dependencies
go mod tidy

# Build the application
go build -o typeracer-tui .

# Run tests
go test ./...
```

### Testing

```bash
# Test practice mode
./typeracer-tui -mode practice

# Test server mode (in one terminal)
./typeracer-tui -mode server

# Connect to server (in another terminal)
ssh localhost -p 2222
```

## Configuration

### Server Settings

- **Port**: SSH server port (default: 2222)
- **Max Players**: Maximum players per room (default: 4)
- **Host Key**: Automatically generated if not present

### Quote API

The application fetches random quotes from [quotable.io](https://quotable.io/) API. If the API is unavailable, it falls back to hardcoded quotes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Acknowledgments

- [Charmbracelet](https://charm.sh/) for the amazing TUI libraries
- [Quotable.io](https://quotable.io/) for providing the quote API
- [TypeRacer](https://typeracer.com/) for inspiration
