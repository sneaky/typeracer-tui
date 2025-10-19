package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"typeracer-tui/game"
	"typeracer-tui/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
)

// SSHServer represents the SSH server for multiplayer games
type SSHServer struct {
	manager *game.Manager
	port    string
}

// NewSSHServer creates a new SSH server
func NewSSHServer(port string) *SSHServer {
	return &SSHServer{
		manager: game.NewManager(),
		port:    port,
	}
}

// Start starts the SSH server
func (s *SSHServer) Start() error {
	// Create SSH server
	server, err := wish.NewServer(
		wish.WithAddress(":"+s.port),
		wish.WithHostKeyPath(".ssh/host_key"),
		wish.WithMiddleware(
			s.gameMiddleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create SSH server: %w", err)
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting TypeRacer SSH server on port %s...", s.port)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("SSH server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown SSH server: %w", err)
	}

	return nil
}

// gameMiddleware creates middleware for handling game sessions
func (s *SSHServer) gameMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(session ssh.Session) {
			// Get player info from SSH session
			playerID := session.User()
			playerName := session.User() // Use username as display name

			// Add player to manager
			_, err := s.manager.AddPlayer(playerID, playerName)
			if err != nil {
				log.Printf("Failed to add player %s: %v", playerID, err)
				session.Close()
				return
			}

			log.Printf("Player %s (%s) connected", playerName, playerID)

			// Try to find an available lobby or create one
			lobby := s.findOrCreateLobby()
			if lobby == nil {
				log.Printf("Failed to create lobby for player %s", playerID)
				session.Close()
				return
			}

			// Join lobby
			if err := s.manager.JoinLobby(playerID, lobby.ID); err != nil {
				log.Printf("Failed to join lobby for player %s: %v", playerID, err)
				session.Close()
				return
			}

			// Create Bubble Tea program
			model := ui.NewLobbyModel(s.manager, playerID, playerName, lobby.ID, lobby.MaxPlayers)
			program := tea.NewProgram(model, tea.WithAltScreen())

			// Handle lobby updates and game transitions
			go s.handlePlayerUpdates(playerID, program)

			// Start the program
			if err := program.Start(); err != nil {
				log.Printf("Error starting program for player %s: %v", playerID, err)
			}

			// Cleanup on disconnect
			s.manager.RemovePlayer(playerID)
			log.Printf("Player %s disconnected", playerID)
			session.Close()
		}
	}
}

// findOrCreateLobby finds an available lobby or creates a new one
func (s *SSHServer) findOrCreateLobby() *game.Lobby {
	// Try to find an available lobby
	availableLobbies := s.manager.GetAvailableLobbies()
	for _, lobby := range availableLobbies {
		if len(lobby.GetPlayers()) < lobby.MaxPlayers {
			return lobby
		}
	}

	// Create new lobby with default settings
	lobby, err := s.manager.CreateLobby(4) // Default to 4 players max
	if err != nil {
		log.Printf("Failed to create lobby: %v", err)
		return nil
	}

	return lobby
}

// handlePlayerUpdates handles player updates and game transitions
func (s *SSHServer) handlePlayerUpdates(playerID string, program *tea.Program) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Check if player is in a lobby that's ready to start
		if lobby, exists := s.manager.GetLobby(playerID); exists {
			if lobby.IsReady() {
				// Start session from lobby
				session, err := s.manager.StartSessionFromLobby(lobby.ID)
				if err != nil {
					log.Printf("Failed to start session from lobby: %v", err)
					continue
				}

				// Send start game message
				program.Send(ui.StartGameMsg{SessionID: session.ID})
				return
			}
		}

		// Check if player is in an active session
		if session, exists := s.manager.GetSession(playerID); exists {
			// Session is active, send refresh message
			program.Send(ui.RefreshGameMsg{})

			// Check if session is finished
			if session.IsFinished {
				return
			}
		}
	}
}

// generateHostKey generates a host key if it doesn't exist
func generateHostKey() error {
	keyPath := ".ssh/host_key"
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		// Create .ssh directory if it doesn't exist
		if err := os.MkdirAll(".ssh", 0700); err != nil {
			return fmt.Errorf("failed to create .ssh directory: %w", err)
		}

		// Generate host key (this is a simplified approach)
		// In production, you'd want to use proper SSH key generation
		log.Printf("Host key not found at %s. Please generate one manually:", keyPath)
		log.Printf("ssh-keygen -t ed25519 -f %s -N ''", keyPath)
		return fmt.Errorf("host key not found")
	}
	return nil
}
