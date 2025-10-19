package main

import (
	"flag"
	"fmt"
	"log"

	"typeracer-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	var (
		mode    = flag.String("mode", "practice", "Mode: 'practice' or 'server'")
		port    = flag.String("port", "2222", "SSH server port (server mode only)")
		players = flag.Int("players", 4, "Maximum players per room (server mode only)")
		help    = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	switch *mode {
	case "practice":
		runPracticeMode()
	case "server":
		runServerMode(*port, *players)
	default:
		log.Fatalf("Invalid mode: %s. Use 'practice' or 'server'", *mode)
	}
}

// runPracticeMode runs the single-player practice mode
func runPracticeMode() {
	fmt.Println("Starting TypeRacer Practice Mode...")

	model := ui.NewPracticeModel()
	program := tea.NewProgram(model, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		log.Fatalf("Error running practice mode: %v", err)
	}
}

// runServerMode runs the SSH server for multiplayer games
func runServerMode(port string, maxPlayers int) {
	fmt.Printf("Starting TypeRacer Server on port %s (max %d players per room)...\n", port, maxPlayers)

	server := NewSSHServer(port)

	// Check for host key
	if err := generateHostKey(); err != nil {
		log.Printf("Warning: %v", err)
		log.Println("You can generate a host key with:")
		log.Printf("ssh-keygen -t ed25519 -f .ssh/host_key -N ''")
		log.Println("Or the server will create one automatically.")
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}

// showHelp displays help information
func showHelp() {
	fmt.Println("TypeRacer TUI - Terminal-based typing race game")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  typeracer-tui [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -mode string")
	fmt.Println("        Mode to run: 'practice' or 'server' (default: practice)")
	fmt.Println("  -port string")
	fmt.Println("        SSH server port for server mode (default: 2222)")
	fmt.Println("  -players int")
	fmt.Println("        Maximum players per room for server mode (default: 4)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Run practice mode")
	fmt.Println("  typeracer-tui")
	fmt.Println("  typeracer-tui -mode practice")
	fmt.Println()
	fmt.Println("  # Run server mode")
	fmt.Println("  typeracer-tui -mode server")
	fmt.Println("  typeracer-tui -mode server -port 2222 -players 4")
	fmt.Println()
	fmt.Println("  # Connect to server")
	fmt.Println("  ssh localhost -p 2222")
	fmt.Println()
	fmt.Println("Practice Mode:")
	fmt.Println("  - Single-player typing practice")
	fmt.Println("  - Real-time WPM and accuracy tracking")
	fmt.Println("  - Visual feedback for correct/incorrect typing")
	fmt.Println("  - No network connection required")
	fmt.Println()
	fmt.Println("Server Mode:")
	fmt.Println("  - Multiplayer typing races over SSH")
	fmt.Println("  - Lobby system for player matchmaking")
	fmt.Println("  - Real-time opponent progress tracking")
	fmt.Println("  - Configurable room sizes (2-4 players)")
	fmt.Println("  - 3-2-1-GO countdown before races")
	fmt.Println()
	fmt.Println("Controls:")
	fmt.Println("  - Type the displayed text as fast and accurately as possible")
	fmt.Println("  - Backspace to correct mistakes")
	fmt.Println("  - Ctrl+C or Esc to quit")
	fmt.Println("  - 'r' to restart (practice mode)")
	fmt.Println("  - 'q' to quit (results screen)")
}
