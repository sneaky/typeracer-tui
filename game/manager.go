package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"typeracer-tui/quotes"

	"github.com/google/uuid"
)

// Manager handles game sessions and player matchmaking
type Manager struct {
	sessions     map[string]*Session
	players      map[string]*Player
	lobbies      map[string]*Lobby
	mu           sync.RWMutex
	quoteFetcher *quotes.Fetcher
}

// Lobby represents a waiting area for players
type Lobby struct {
	ID         string             `json:"id"`
	Players    map[string]*Player `json:"players"`
	MaxPlayers int                `json:"max_players"`
	CreatedAt  time.Time          `json:"created_at"`
	mu         sync.RWMutex
}

// NewManager creates a new game manager
func NewManager() *Manager {
	return &Manager{
		sessions:     make(map[string]*Session),
		players:      make(map[string]*Player),
		lobbies:      make(map[string]*Lobby),
		quoteFetcher: quotes.NewFetcher(),
	}
}

// AddPlayer adds a player to the system
func (m *Manager) AddPlayer(playerID, playerName string) (*Player, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if player already exists
	if _, exists := m.players[playerID]; exists {
		return nil, fmt.Errorf("player already exists")
	}

	player := NewPlayer(playerID, playerName, "")
	m.players[playerID] = player

	log.Printf("Player %s (%s) added to system", playerName, playerID)
	return player, nil
}

// RemovePlayer removes a player from the system
func (m *Manager) RemovePlayer(playerID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove from all sessions
	for sessionID, session := range m.sessions {
		session.RemovePlayer(playerID)
		if len(session.GetPlayers()) == 0 {
			delete(m.sessions, sessionID)
		}
	}

	// Remove from all lobbies
	for lobbyID, lobby := range m.lobbies {
		lobby.RemovePlayer(playerID)
		if len(lobby.GetPlayers()) == 0 {
			delete(m.lobbies, lobbyID)
		}
	}

	delete(m.players, playerID)
	log.Printf("Player %s removed from system", playerID)
}

// GetPlayer returns a player by ID
func (m *Manager) GetPlayer(playerID string) (*Player, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	player, exists := m.players[playerID]
	return player, exists
}

// CreateLobby creates a new lobby for players to wait
func (m *Manager) CreateLobby(maxPlayers int) (*Lobby, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lobbyID := uuid.New().String()
	lobby := &Lobby{
		ID:         lobbyID,
		Players:    make(map[string]*Player),
		MaxPlayers: maxPlayers,
		CreatedAt:  time.Now(),
	}

	m.lobbies[lobbyID] = lobby
	log.Printf("Created lobby %s with max %d players", lobbyID, maxPlayers)
	return lobby, nil
}

// JoinLobby adds a player to a lobby
func (m *Manager) JoinLobby(playerID, lobbyID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	player, exists := m.players[playerID]
	if !exists {
		return fmt.Errorf("player not found")
	}

	lobby, exists := m.lobbies[lobbyID]
	if !exists {
		return fmt.Errorf("lobby not found")
	}

	if len(lobby.Players) >= lobby.MaxPlayers {
		return fmt.Errorf("lobby is full")
	}

	lobby.Players[playerID] = player
	player.SessionID = lobbyID

	log.Printf("Player %s joined lobby %s", playerID, lobbyID)
	return nil
}

// LeaveLobby removes a player from a lobby
func (m *Manager) LeaveLobby(playerID, lobbyID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if lobby, exists := m.lobbies[lobbyID]; exists {
		lobby.RemovePlayer(playerID)
		if len(lobby.GetPlayers()) == 0 {
			delete(m.lobbies, lobbyID)
		}
	}
}

// StartSessionFromLobby starts a session from a lobby
func (m *Manager) StartSessionFromLobby(lobbyID string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	lobby, exists := m.lobbies[lobbyID]
	if !exists {
		return nil, fmt.Errorf("lobby not found")
	}

	if len(lobby.Players) < 2 {
		return nil, fmt.Errorf("not enough players to start session")
	}

	// Fetch a random quote
	quote := m.quoteFetcher.FetchRandomQuoteWithFallback()

	// Create session
	sessionID := uuid.New().String()
	session := NewSession(sessionID, quote.Content, quote.Author, lobby.MaxPlayers)

	// Add all players from lobby to session
	for _, player := range lobby.Players {
		player.SessionID = sessionID
		session.AddPlayer(player)
	}

	m.sessions[sessionID] = session

	// Remove lobby
	delete(m.lobbies, lobbyID)

	log.Printf("Started session %s with %d players", sessionID, len(session.Players))
	return session, nil
}

// GetSession returns a session by ID
func (m *Manager) GetSession(sessionID string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[sessionID]
	return session, exists
}

// GetLobby returns a lobby by ID
func (m *Manager) GetLobby(lobbyID string) (*Lobby, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lobby, exists := m.lobbies[lobbyID]
	return lobby, exists
}

// GetAvailableLobbies returns all lobbies that can accept more players
func (m *Manager) GetAvailableLobbies() []*Lobby {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var available []*Lobby
	for _, lobby := range m.lobbies {
		if len(lobby.Players) < lobby.MaxPlayers {
			available = append(available, lobby)
		}
	}
	return available
}

// UpdatePlayerProgress updates a player's progress in their session
func (m *Manager) UpdatePlayerProgress(playerID, typedInput string) error {
	m.mu.RLock()
	session, exists := m.sessions[playerID]
	m.mu.RUnlock()

	if !exists {
		// Try to find session by player
		for _, s := range m.sessions {
			if _, playerExists := s.GetPlayer(playerID); playerExists {
				session = s
				break
			}
		}
	}

	if session == nil {
		return fmt.Errorf("player not in any session")
	}

	session.UpdatePlayerProgress(playerID, typedInput)
	return nil
}

// GetSystemStatus returns the current system status
func (m *Manager) GetSystemStatus() SystemStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return SystemStatus{
		TotalPlayers:   len(m.players),
		ActiveSessions: len(m.sessions),
		ActiveLobbies:  len(m.lobbies),
	}
}

// SystemStatus represents the current system status
type SystemStatus struct {
	TotalPlayers   int `json:"total_players"`
	ActiveSessions int `json:"active_sessions"`
	ActiveLobbies  int `json:"active_lobbies"`
}

// Lobby methods
func (l *Lobby) AddPlayer(player *Player) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.Players) >= l.MaxPlayers {
		return fmt.Errorf("lobby is full")
	}

	l.Players[player.ID] = player
	return nil
}

func (l *Lobby) RemovePlayer(playerID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.Players, playerID)
}

func (l *Lobby) GetPlayers() []*Player {
	l.mu.RLock()
	defer l.mu.RUnlock()

	players := make([]*Player, 0, len(l.Players))
	for _, player := range l.Players {
		players = append(players, player)
	}
	return players
}

func (l *Lobby) IsReady() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return len(l.Players) >= 2
}
