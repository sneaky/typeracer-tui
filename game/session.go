package game

import (
	"fmt"
	"sync"
	"time"
)

// Session represents a game session
type Session struct {
	ID         string             `json:"id"`
	Prompt     string             `json:"prompt"`
	Author     string             `json:"author"`
	Players    map[string]*Player `json:"players"`
	MaxPlayers int                `json:"max_players"`
	StartTime  time.Time          `json:"start_time"`
	EndTime    time.Time          `json:"end_time"`
	IsActive   bool               `json:"is_active"`
	IsFinished bool               `json:"is_finished"`
	Countdown  int                `json:"countdown"`
	mu         sync.RWMutex
}

// NewSession creates a new game session
func NewSession(id, prompt, author string, maxPlayers int) *Session {
	return &Session{
		ID:         id,
		Prompt:     prompt,
		Author:     author,
		Players:    make(map[string]*Player),
		MaxPlayers: maxPlayers,
		IsActive:   false,
		IsFinished: false,
		Countdown:  0,
	}
}

// AddPlayer adds a player to the session
func (s *Session) AddPlayer(player *Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Players) >= s.MaxPlayers {
		return fmt.Errorf("session is full")
	}

	if s.IsActive {
		return fmt.Errorf("session has already started")
	}

	s.Players[player.ID] = player
	return nil
}

// RemovePlayer removes a player from the session
func (s *Session) RemovePlayer(playerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Players, playerID)
}

// GetPlayer returns a player by ID
func (s *Session) GetPlayer(playerID string) (*Player, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, exists := s.Players[playerID]
	return player, exists
}

// GetPlayers returns all players in the session
func (s *Session) GetPlayers() []*Player {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]*Player, 0, len(s.Players))
	for _, player := range s.Players {
		players = append(players, player)
	}
	return players
}

// IsReady checks if the session is ready to start
func (s *Session) IsReady() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.Players) >= 2 && !s.IsActive
}

// Start begins the session with a countdown
func (s *Session) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IsActive {
		return fmt.Errorf("session is already active")
	}

	if len(s.Players) < 2 {
		return fmt.Errorf("not enough players to start")
	}

	s.IsActive = true
	s.StartTime = time.Now()
	s.Countdown = 3

	// Start countdown in a goroutine
	go s.runCountdown()

	return nil
}

// runCountdown runs the 3-2-1-GO countdown
func (s *Session) runCountdown() {
	for i := 3; i > 0; i-- {
		s.mu.Lock()
		s.Countdown = i
		s.mu.Unlock()

		time.Sleep(1 * time.Second)
	}

	s.mu.Lock()
	s.Countdown = 0
	s.mu.Unlock()
}

// UpdatePlayerProgress updates a player's progress
func (s *Session) UpdatePlayerProgress(playerID, typedInput string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if player, exists := s.Players[playerID]; exists {
		player.UpdateProgress(typedInput, s.Prompt)

		// Check if player finished
		if player.IsComplete(len(s.Prompt)) && !player.IsFinished {
			player.Finish()
		}

		// Check if all players finished
		s.checkCompletion()
	}
}

// checkCompletion checks if all players have finished
func (s *Session) checkCompletion() {
	if s.IsFinished {
		return
	}

	allFinished := true
	for _, player := range s.Players {
		if !player.IsFinished {
			allFinished = false
			break
		}
	}

	if allFinished {
		s.IsFinished = true
		s.EndTime = time.Now()
	}
}

// GetLeaderboard returns players sorted by completion time
func (s *Session) GetLeaderboard() []*Player {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make([]*Player, 0, len(s.Players))
	for _, player := range s.Players {
		players = append(players, player)
	}

	// Sort by finish time (earliest first)
	for i := 0; i < len(players); i++ {
		for j := i + 1; j < len(players); j++ {
			if players[i].IsFinished && players[j].IsFinished {
				if players[i].EndTime.After(players[j].EndTime) {
					players[i], players[j] = players[j], players[i]
				}
			} else if players[i].IsFinished && !players[j].IsFinished {
				// Finished players come first
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	return players
}

// GetStatus returns the current session status
func (s *Session) GetStatus() SessionStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return SessionStatus{
		ID:          s.ID,
		PlayerCount: len(s.Players),
		MaxPlayers:  s.MaxPlayers,
		IsActive:    s.IsActive,
		IsFinished:  s.IsFinished,
		Countdown:   s.Countdown,
	}
}

// SessionStatus represents the current status of a session
type SessionStatus struct {
	ID          string `json:"id"`
	PlayerCount int    `json:"player_count"`
	MaxPlayers  int    `json:"max_players"`
	IsActive    bool   `json:"is_active"`
	IsFinished  bool   `json:"is_finished"`
	Countdown   int    `json:"countdown"`
}

