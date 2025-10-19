package game

import (
	"time"
)

// Player represents a player in the game
type Player struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SessionID    string    `json:"session_id"`
	CurrentPos   int       `json:"current_pos"`
	TypedInput   string    `json:"typed_input"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsFinished   bool      `json:"is_finished"`
	WPM          float64   `json:"wpm"`
	Accuracy     float64   `json:"accuracy"`
	CorrectChars int       `json:"correct_chars"`
	TotalChars   int       `json:"total_chars"`
	LastUpdate   time.Time `json:"last_update"`
}

// NewPlayer creates a new player
func NewPlayer(id, name, sessionID string) *Player {
	return &Player{
		ID:         id,
		Name:       name,
		SessionID:  sessionID,
		CurrentPos: 0,
		TypedInput: "",
		StartTime:  time.Now(),
		IsFinished: false,
		WPM:        0.0,
		Accuracy:   0.0,
		LastUpdate: time.Now(),
	}
}

// UpdateProgress updates the player's typing progress
func (p *Player) UpdateProgress(typedInput string, prompt string) {
	p.TypedInput = typedInput
	p.CurrentPos = len(typedInput)
	p.LastUpdate = time.Now()

	// Calculate accuracy
	p.calculateAccuracy(prompt)

	// Calculate WPM
	p.calculateWPM()
}

// calculateAccuracy calculates the player's typing accuracy
func (p *Player) calculateAccuracy(prompt string) {
	if len(prompt) == 0 {
		p.Accuracy = 0.0
		return
	}

	correct := 0
	total := len(p.TypedInput)
	if total > len(prompt) {
		total = len(prompt)
	}

	for i := 0; i < total; i++ {
		if i < len(p.TypedInput) && i < len(prompt) && p.TypedInput[i] == prompt[i] {
			correct++
		}
	}

	p.CorrectChars = correct
	p.TotalChars = total

	if total > 0 {
		p.Accuracy = float64(correct) / float64(total) * 100.0
	} else {
		p.Accuracy = 0.0
	}
}

// calculateWPM calculates words per minute
func (p *Player) calculateWPM() {
	if p.IsFinished {
		// Use end time for final WPM calculation
		elapsed := p.EndTime.Sub(p.StartTime).Minutes()
		if elapsed > 0 {
			// WPM = (correct characters / 5) / minutes
			p.WPM = float64(p.CorrectChars) / 5.0 / elapsed
		}
	} else {
		// Use current time for live WPM calculation
		elapsed := time.Since(p.StartTime).Minutes()
		if elapsed > 0 {
			p.WPM = float64(p.CorrectChars) / 5.0 / elapsed
		}
	}
}

// Finish marks the player as finished and calculates final stats
func (p *Player) Finish() {
	p.IsFinished = true
	p.EndTime = time.Now()
	p.calculateWPM()
}

// GetProgress returns the progress percentage (0-100)
func (p *Player) GetProgress(promptLength int) float64 {
	if promptLength == 0 {
		return 0.0
	}
	return float64(p.CurrentPos) / float64(promptLength) * 100.0
}

// IsComplete checks if the player has completed the typing
func (p *Player) IsComplete(promptLength int) bool {
	return p.CurrentPos >= promptLength
}

