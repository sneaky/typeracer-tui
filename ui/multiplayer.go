package ui

import (
	"fmt"
	"strings"
	"time"

	"typeracer-tui/game"

	tea "github.com/charmbracelet/bubbletea"
)

// MultiplayerModel represents the multiplayer game mode
type MultiplayerModel struct {
	manager       *game.Manager
	playerID      string
	playerName    string
	sessionID     string
	session       *game.Session
	typedInput    string
	startTime     time.Time
	isFinished    bool
	wpm           float64
	accuracy      float64
	correctChars  int
	width         int
	height        int
	showResults   bool
	refreshTicker *time.Ticker
}

// NewMultiplayerModel creates a new multiplayer model
func NewMultiplayerModel(manager *game.Manager, playerID, playerName, sessionID string) *MultiplayerModel {
	return &MultiplayerModel{
		manager:    manager,
		playerID:   playerID,
		playerName: playerName,
		sessionID:  sessionID,
		width:      80,
		height:     24,
	}
}

// Init initializes the multiplayer model
func (m *MultiplayerModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.startRefreshTicker(),
	)
}

// startRefreshTicker starts a ticker to refresh game state
func (m *MultiplayerModel) startRefreshTicker() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond)
		return RefreshGameMsg{}
	}
}

// Update handles messages and updates the model
func (m *MultiplayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.showResults {
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "backspace":
				if len(m.typedInput) > 0 {
					m.typedInput = m.typedInput[:len(m.typedInput)-1]
					m.updateProgress()
				}
			default:
				if len(msg.String()) == 1 {
					m.typedInput += msg.String()
					m.updateProgress()

					// Check if finished
					if m.isComplete() && !m.isFinished {
						m.finish()
					}
				}
			}
		}
		return m, nil

	case RefreshGameMsg:
		// Update session state
		if session, exists := m.manager.GetSession(m.sessionID); exists {
			m.session = session

			// Check if game has started
			if session.IsActive && session.Countdown == 0 && m.startTime.IsZero() {
				m.startTime = time.Now()
			}

			// Check if game is finished
			if session.IsFinished && !m.showResults {
				m.showResults = true
			}
		}
		return m, m.startRefreshTicker()
	}

	return m, nil
}

// View renders the multiplayer game UI
func (m *MultiplayerModel) View() string {
	if m.session == nil {
		return TitleStyle.Render("Loading game...")
	}

	if m.showResults {
		return m.renderResults()
	}

	if m.session.Countdown > 0 {
		return m.renderCountdown()
	}

	return m.renderGame()
}

// renderCountdown renders the countdown screen
func (m *MultiplayerModel) renderCountdown() string {
	var content strings.Builder

	// Title
	content.WriteString(TitleStyle.Render("Get Ready!"))
	content.WriteString("\n\n")

	// Countdown
	countdownText := fmt.Sprintf("%d", m.session.Countdown)
	content.WriteString(CountdownStyle.Render(countdownText))
	content.WriteString("\n\n")

	// Players list
	content.WriteString(m.renderPlayersList())
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Get ready to type!"))

	return content.String()
}

// renderGame renders the main game interface
func (m *MultiplayerModel) renderGame() string {
	var content strings.Builder

	// Title
	content.WriteString(TitleStyle.Render("TypeRacer Multiplayer"))
	content.WriteString("\n\n")

	// Quote author
	if m.session.Author != "" {
		content.WriteString(SubtitleStyle.Render(fmt.Sprintf("— %s", m.session.Author)))
		content.WriteString("\n\n")
	}

	// Typing area
	typingBox := MainBoxStyle.Width(m.width - 4).Render(
		StyleTypingText(m.session.Prompt, m.typedInput),
	)
	content.WriteString(typingBox)
	content.WriteString("\n\n")

	// Stats
	stats := m.renderStats()
	content.WriteString(stats)
	content.WriteString("\n\n")

	// Progress bar
	progress := CreateProgressBar(len(m.typedInput), len(m.session.Prompt), m.width-10)
	content.WriteString(ProgressBoxStyle.Render(progress))
	content.WriteString("\n\n")

	// Opponents
	content.WriteString(m.renderOpponents())
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Type as fast and accurately as possible!"))

	return content.String()
}

// renderResults renders the completion screen
func (m *MultiplayerModel) renderResults() string {
	var content strings.Builder

	// Title
	content.WriteString(LeaderboardTitleStyle.Render("Race Complete!"))
	content.WriteString("\n\n")

	// Leaderboard
	content.WriteString(m.renderLeaderboard())
	content.WriteString("\n\n")

	// Your results
	content.WriteString(m.renderYourResults())
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Press 'q' to quit"))

	return content.String()
}

// renderStats renders the current stats
func (m *MultiplayerModel) renderStats() string {
	var stats strings.Builder

	// WPM
	stats.WriteString(WPMTextStyle.Render(fmt.Sprintf("WPM: %s", FormatWPM(m.wpm))))
	stats.WriteString("  ")

	// Accuracy
	stats.WriteString(AccuracyTextStyle.Render(fmt.Sprintf("Accuracy: %s", FormatAccuracy(m.accuracy))))
	stats.WriteString("  ")

	// Time
	elapsed := time.Since(m.startTime).Seconds()
	stats.WriteString(TimeTextStyle.Render(fmt.Sprintf("Time: %s", FormatDuration(elapsed))))

	return StatsBoxStyle.Render(stats.String())
}

// renderPlayersList renders the list of players
func (m *MultiplayerModel) renderPlayersList() string {
	var content strings.Builder

	players := m.session.GetPlayers()
	content.WriteString(PlayerNameStyle.Render(fmt.Sprintf("Players (%d)", len(players))))
	content.WriteString("\n")

	for i, player := range players {
		playerText := fmt.Sprintf("%d. %s", i+1, player.Name)
		if player.ID == m.playerID {
			playerText += " (You)"
		}
		content.WriteString(PlayerNameStyle.Render(playerText))
		content.WriteString("\n")
	}

	return MainBoxStyle.Width(m.width - 4).Render(content.String())
}

// renderOpponents renders opponent progress
func (m *MultiplayerModel) renderOpponents() string {
	var content strings.Builder

	content.WriteString(PlayerNameStyle.Render("Opponents"))
	content.WriteString("\n")

	players := m.session.GetPlayers()
	for i, player := range players {
		if player.ID == m.playerID {
			continue // Skip self
		}

		// Player name and racer
		racer := CreateRacerIndicator(i, player.IsFinished)
		playerText := fmt.Sprintf("%s %s", racer, player.Name)
		if player.IsFinished {
			playerText += " ✓"
		}
		content.WriteString(PlayerNameStyle.Render(playerText))
		content.WriteString("\n")

		// Progress bar
		progress := CreateProgressBar(player.CurrentPos, len(m.session.Prompt), 30)
		content.WriteString(ProgressBoxStyle.Render(progress))
		content.WriteString("\n")

		// WPM
		content.WriteString(PlayerWPMStyle.Render(fmt.Sprintf("WPM: %s", FormatWPM(player.WPM))))
		content.WriteString("\n\n")
	}

	return MainBoxStyle.Width(m.width - 4).Render(content.String())
}

// renderLeaderboard renders the final leaderboard
func (m *MultiplayerModel) renderLeaderboard() string {
	var content strings.Builder

	leaderboard := m.session.GetLeaderboard()
	content.WriteString(LeaderboardTitleStyle.Render("Final Results"))
	content.WriteString("\n\n")

	for i, player := range leaderboard {
		position := i + 1
		positionText := fmt.Sprintf("%d.", position)

		// Position styling
		if position == 1 {
			positionText = "🥇 " + positionText
		} else if position == 2 {
			positionText = "🥈 " + positionText
		} else if position == 3 {
			positionText = "🥉 " + positionText
		}

		// Player info
		playerInfo := fmt.Sprintf("%s %s", positionText, player.Name)
		if player.ID == m.playerID {
			playerInfo += " (You)"
		}

		content.WriteString(LeaderboardEntryStyle.Render(playerInfo))
		content.WriteString("\n")

		// Stats
		stats := fmt.Sprintf("WPM: %s | Accuracy: %s",
			FormatWPM(player.WPM),
			FormatAccuracy(player.Accuracy))
		content.WriteString(LeaderboardWPMStyle.Render(stats))
		content.WriteString("\n\n")
	}

	return MainBoxStyle.Width(m.width - 4).Render(content.String())
}

// renderYourResults renders your personal results
func (m *MultiplayerModel) renderYourResults() string {
	results := fmt.Sprintf(
		"Your Results:\nWPM: %s | Accuracy: %s | Time: %s",
		FormatWPM(m.wpm),
		FormatAccuracy(m.accuracy),
		FormatDuration(time.Since(m.startTime).Seconds()),
	)

	return MainBoxStyle.Width(m.width - 4).Render(results)
}

// updateProgress updates the player's progress
func (m *MultiplayerModel) updateProgress() {
	if m.session == nil {
		return
	}

	// Update local stats
	m.calculateStats()

	// Update in game manager
	m.manager.UpdatePlayerProgress(m.playerID, m.typedInput)
}

// calculateStats calculates local statistics
func (m *MultiplayerModel) calculateStats() {
	if m.session == nil {
		return
	}

	// Calculate accuracy
	correct := 0
	total := len(m.typedInput)
	if total > len(m.session.Prompt) {
		total = len(m.session.Prompt)
	}

	for i := 0; i < total; i++ {
		if i < len(m.typedInput) && i < len(m.session.Prompt) && m.typedInput[i] == m.session.Prompt[i] {
			correct++
		}
	}

	m.correctChars = correct

	if total > 0 {
		m.accuracy = float64(correct) / float64(total) * 100.0
	} else {
		m.accuracy = 0.0
	}

	// Calculate WPM
	elapsed := time.Since(m.startTime).Minutes()
	if elapsed > 0 {
		m.wpm = float64(correct) / 5.0 / elapsed
	}
}

// isComplete checks if the typing is complete
func (m *MultiplayerModel) isComplete() bool {
	return len(m.typedInput) >= len(m.session.Prompt)
}

// finish marks the player as finished
func (m *MultiplayerModel) finish() {
	m.isFinished = true
}

// RefreshGameMsg represents a message to refresh game state
type RefreshGameMsg struct{}
