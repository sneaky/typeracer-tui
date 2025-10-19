package ui

import (
	"fmt"
	"strings"
	"time"

	"typeracer-tui/quotes"

	tea "github.com/charmbracelet/bubbletea"
)

// PracticeModel represents the single-player practice mode
type PracticeModel struct {
	quote        *quotes.Quote
	typedInput   string
	startTime    time.Time
	endTime      time.Time
	isFinished   bool
	wpm          float64
	accuracy     float64
	correctChars int
	totalChars   int
	width        int
	height       int
	showResults  bool
}

// NewPracticeModel creates a new practice mode model
func NewPracticeModel() *PracticeModel {
	return &PracticeModel{
		width:  80,
		height: 24,
	}
}

// Init initializes the practice model
func (m *PracticeModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchQuote(),
		tea.EnterAltScreen,
	)
}

// fetchQuote fetches a random quote
func (m *PracticeModel) fetchQuote() tea.Cmd {
	return func() tea.Msg {
		fetcher := quotes.NewFetcher()
		quote := fetcher.FetchRandomQuoteWithFallback()
		return QuoteMsg{Quote: quote}
	}
}

// Update handles messages and updates the model
func (m *PracticeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case "r", "enter":
				// Restart practice
				newModel := NewPracticeModel()
				newModel.width = m.width
				newModel.height = m.height
				return newModel, newModel.fetchQuote()
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "backspace":
				if len(m.typedInput) > 0 {
					m.typedInput = m.typedInput[:len(m.typedInput)-1]
					m.updateStats()
				}
			default:
				if len(msg.String()) == 1 {
					m.typedInput += msg.String()
					m.updateStats()

					// Check if finished
					if m.isComplete() && !m.isFinished {
						m.finish()
					}
				}
			}
		}
		return m, nil

	case QuoteMsg:
		m.quote = msg.Quote
		m.startTime = time.Now()
		return m, nil
	}

	return m, nil
}

// View renders the practice mode UI
func (m *PracticeModel) View() string {
	if m.quote == nil {
		return TitleStyle.Render("Loading quote...")
	}

	if m.showResults {
		return m.renderResults()
	}

	return m.renderPractice()
}

// renderPractice renders the practice interface
func (m *PracticeModel) renderPractice() string {
	var content strings.Builder

	// Title
	content.WriteString(TitleStyle.Render("TypeRacer Practice"))
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Type the text below as fast and accurately as possible"))
	content.WriteString("\n\n")

	// Quote author
	if m.quote.Author != "" {
		content.WriteString(SubtitleStyle.Render(fmt.Sprintf("— %s", m.quote.Author)))
		content.WriteString("\n\n")
	}

	// Typing area
	typingBox := MainBoxStyle.Width(m.width - 4).Render(
		StyleTypingText(m.quote.Content, m.typedInput),
	)
	content.WriteString(typingBox)
	content.WriteString("\n\n")

	// Stats
	stats := m.renderStats()
	content.WriteString(stats)
	content.WriteString("\n\n")

	// Progress bar
	progress := CreateProgressBar(len(m.typedInput), len(m.quote.Content), m.width-10)
	content.WriteString(ProgressBoxStyle.Render(progress))
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Press Ctrl+C or Esc to quit"))

	return content.String()
}

// renderResults renders the completion screen
func (m *PracticeModel) renderResults() string {
	var content strings.Builder

	// Title
	content.WriteString(LeaderboardTitleStyle.Render("Practice Complete!"))
	content.WriteString("\n\n")

	// Results box
	results := fmt.Sprintf(
		"Words Per Minute: %s\nAccuracy: %s\nTime: %s\nCharacters: %d/%d",
		FormatWPM(m.wpm),
		FormatAccuracy(m.accuracy),
		FormatDuration(m.endTime.Sub(m.startTime).Seconds()),
		m.correctChars,
		len(m.quote.Content),
	)

	resultsBox := MainBoxStyle.Width(m.width - 4).Render(results)
	content.WriteString(resultsBox)
	content.WriteString("\n\n")

	// Instructions
	content.WriteString(InstructionStyle.Render("Press 'r' to restart or 'q' to quit"))

	return content.String()
}

// renderStats renders the current stats
func (m *PracticeModel) renderStats() string {
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

// updateStats updates the player's statistics
func (m *PracticeModel) updateStats() {
	if m.quote == nil {
		return
	}

	// Calculate accuracy
	m.calculateAccuracy()

	// Calculate WPM
	m.calculateWPM()
}

// calculateAccuracy calculates typing accuracy
func (m *PracticeModel) calculateAccuracy() {
	if len(m.quote.Content) == 0 {
		m.accuracy = 0.0
		return
	}

	correct := 0
	total := len(m.typedInput)
	if total > len(m.quote.Content) {
		total = len(m.quote.Content)
	}

	for i := 0; i < total; i++ {
		if i < len(m.typedInput) && i < len(m.quote.Content) && m.typedInput[i] == m.quote.Content[i] {
			correct++
		}
	}

	m.correctChars = correct
	m.totalChars = total

	if total > 0 {
		m.accuracy = float64(correct) / float64(total) * 100.0
	} else {
		m.accuracy = 0.0
	}
}

// calculateWPM calculates words per minute
func (m *PracticeModel) calculateWPM() {
	elapsed := time.Since(m.startTime).Minutes()
	if elapsed > 0 {
		// WPM = (correct characters / 5) / minutes
		m.wpm = float64(m.correctChars) / 5.0 / elapsed
	}
}

// isComplete checks if the typing is complete
func (m *PracticeModel) isComplete() bool {
	return len(m.typedInput) >= len(m.quote.Content)
}

// finish marks the practice as finished
func (m *PracticeModel) finish() {
	m.isFinished = true
	m.endTime = time.Now()
	m.showResults = true
}

// QuoteMsg represents a message containing a quote
type QuoteMsg struct {
	Quote *quotes.Quote
}
