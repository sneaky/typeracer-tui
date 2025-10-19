package ui

import (
	"fmt"
	"strings"
	"time"

	"typeracer-tui/game"

	tea "github.com/charmbracelet/bubbletea"
)

// LobbyModel represents the lobby waiting screen
type LobbyModel struct {
	manager       *game.Manager
	playerID      string
	playerName    string
	lobbyID       string
	players       []*game.Player
	maxPlayers    int
	width         int
	height        int
	refreshTicker *time.Ticker
}

// NewLobbyModel creates a new lobby model
func NewLobbyModel(manager *game.Manager, playerID, playerName, lobbyID string, maxPlayers int) *LobbyModel {
	return &LobbyModel{
		manager:    manager,
		playerID:   playerID,
		playerName: playerName,
		lobbyID:    lobbyID,
		maxPlayers: maxPlayers,
		width:      80,
		height:     24,
	}
}

// Init initializes the lobby model
func (m *LobbyModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.startRefreshTicker(),
	)
}

// startRefreshTicker starts a ticker to refresh lobby state
func (m *LobbyModel) startRefreshTicker() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		return RefreshLobbyMsg{}
	}
}

// Update handles messages and updates the model
func (m *LobbyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "r":
			// Refresh lobby
			return m, m.startRefreshTicker()
		}
		return m, nil

	case RefreshLobbyMsg:
		// Update lobby state
		if lobby, exists := m.manager.GetLobby(m.lobbyID); exists {
			m.players = lobby.GetPlayers()
			m.maxPlayers = lobby.MaxPlayers
		}
		return m, m.startRefreshTicker()

	case StartGameMsg:
		// Game is starting, transition to multiplayer mode
		return NewMultiplayerModel(m.manager, m.playerID, m.playerName, msg.SessionID), nil
	}

	return m, nil
}

// View renders the lobby UI
func (m *LobbyModel) View() string {
	var content strings.Builder

	// Title
	content.WriteString(TitleStyle.Render("TypeRacer Lobby"))
	content.WriteString("\n\n")

	// Lobby info
	lobbyInfo := fmt.Sprintf("Lobby ID: %s", m.lobbyID)
	content.WriteString(SubtitleStyle.Render(lobbyInfo))
	content.WriteString("\n\n")

	// Players list
	content.WriteString(m.renderPlayersList())
	content.WriteString("\n\n")

	// Status
	content.WriteString(m.renderStatus())
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Waiting for players... Press 'r' to refresh, 'q' to quit"))

	return content.String()
}

// renderPlayersList renders the list of connected players
func (m *LobbyModel) renderPlayersList() string {
	var content strings.Builder

	// Players header
	content.WriteString(PlayerNameStyle.Render(fmt.Sprintf("Players (%d/%d)", len(m.players), m.maxPlayers)))
	content.WriteString("\n")

	// Players list
	if len(m.players) == 0 {
		content.WriteString(InstructionStyle.Render("No players connected"))
	} else {
		for i, player := range m.players {
			playerText := fmt.Sprintf("%d. %s", i+1, player.Name)
			if player.ID == m.playerID {
				playerText += " (You)"
			}
			content.WriteString(PlayerNameStyle.Render(playerText))
			content.WriteString("\n")
		}
	}

	return MainBoxStyle.Width(m.width - 4).Render(content.String())
}

// renderStatus renders the current lobby status
func (m *LobbyModel) renderStatus() string {
	var status strings.Builder

	if len(m.players) < 2 {
		status.WriteString(InstructionStyle.Render("Waiting for more players..."))
	} else if len(m.players) < m.maxPlayers {
		status.WriteString(InstructionStyle.Render("Ready to start! Waiting for more players or start now..."))
	} else {
		status.WriteString(SuccessStyle.Render("Lobby is full! Game will start automatically..."))
	}

	return status.String()
}

// RefreshLobbyMsg represents a message to refresh lobby state
type RefreshLobbyMsg struct{}

// StartGameMsg represents a message that the game is starting
type StartGameMsg struct {
	SessionID string
}
