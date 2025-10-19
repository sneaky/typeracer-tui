package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	Green  = lipgloss.Color("#00FF00")
	Red    = lipgloss.Color("#FF0000")
	Yellow = lipgloss.Color("#FFFF00")
	Blue   = lipgloss.Color("#0080FF")
	Purple = lipgloss.Color("#8000FF")
	Orange = lipgloss.Color("#FF8000")

	// Neutral colors
	White    = lipgloss.Color("#FFFFFF")
	Black    = lipgloss.Color("#000000")
	Gray     = lipgloss.Color("#808080")
	DarkGray = lipgloss.Color("#404040")

	// Background colors
	BgDark  = lipgloss.Color("#1a1a1a")
	BgLight = lipgloss.Color("#f0f0f0")
)

// Text styles
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue).
			Margin(1, 0).
			Align(lipgloss.Center)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Margin(0, 0, 1, 0).
			Align(lipgloss.Center)

	// Typing text styles
	CorrectTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Green)

	IncorrectTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Red)

	CurrentTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Yellow).
				Background(DarkGray)

	UntypedTextStyle = lipgloss.NewStyle().
				Foreground(White)

	// Status styles
	WPMTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue)

	AccuracyTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Green)

	TimeTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Orange)

	// Countdown styles
	CountdownStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red).
			Background(White).
			Margin(2, 0).
			Align(lipgloss.Center)

	GoStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Green).
		Background(White).
		Margin(2, 0).
		Align(lipgloss.Center)

	// Box styles
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Blue).
			Padding(1, 2).
			Margin(1, 0)

	StatsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Gray).
			Padding(0, 1).
			Margin(0, 1)

	ProgressBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(Green).
				Padding(0, 1).
				Margin(0, 1)

	// Player list styles
	PlayerNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue)

	PlayerWPMStyle = lipgloss.NewStyle().
			Foreground(Green)

	PlayerProgressStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	// Racer indicators
	RacerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Blue)

	RacerFinishedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Green)

	// Progress bar styles
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(Green).
				Background(DarkGray)

	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(DarkGray).
				Background(DarkGray)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Blue).
			Padding(0, 1).
			Margin(0, 1)

	ButtonActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Green).
				Foreground(Green).
				Padding(0, 1).
				Margin(0, 1)

	// Error and success styles
	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red).
			Background(White).
			Padding(0, 1).
			Margin(1, 0)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Green).
			Background(White).
			Padding(0, 1).
			Margin(1, 0)

	// Instructions
	InstructionStyle = lipgloss.NewStyle().
				Foreground(Gray).
				Italic(true).
				Margin(0, 0, 1, 0)

	// Leaderboard styles
	LeaderboardTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Blue).
				Margin(1, 0).
				Align(lipgloss.Center)

	LeaderboardEntryStyle = lipgloss.NewStyle().
				Margin(0, 0, 0, 1)

	LeaderboardPositionStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(Yellow)

	LeaderboardNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Blue)

	LeaderboardWPMStyle = lipgloss.NewStyle().
				Foreground(Green)

	LeaderboardTimeStyle = lipgloss.NewStyle().
				Foreground(Orange)
)

// Helper functions for styling text with typing progress
func StyleTypingText(prompt, typed string) string {
	if len(typed) == 0 {
		return UntypedTextStyle.Render(prompt)
	}

	var result string
	for i, char := range prompt {
		if i < len(typed) {
			if rune(typed[i]) == char {
				result += CorrectTextStyle.Render(string(char))
			} else {
				result += IncorrectTextStyle.Render(string(char))
			}
		} else if i == len(typed) {
			result += CurrentTextStyle.Render(string(char))
		} else {
			result += UntypedTextStyle.Render(string(char))
		}
	}

	return result
}

// Create a progress bar
func CreateProgressBar(current, total int, width int) string {
	if total == 0 {
		return ProgressBarEmptyStyle.Render(string(make([]rune, width)))
	}

	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return ProgressBarStyle.Render(bar[:filled]) + ProgressBarEmptyStyle.Render(bar[filled:])
}

// Create a racer indicator
func CreateRacerIndicator(position int, isFinished bool) string {
	racers := []string{"🏃", "🚗", "🏎️", "🚴"}
	if position >= len(racers) {
		position = len(racers) - 1
	}

	racer := racers[position]
	if isFinished {
		return RacerFinishedStyle.Render(racer + " ✓")
	}
	return RacerStyle.Render(racer)
}

// Format time duration
func FormatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.1fs", seconds)
	}
	minutes := int(seconds / 60)
	secs := int(seconds) % 60
	return fmt.Sprintf("%dm %ds", minutes, secs)
}

// Format WPM
func FormatWPM(wpm float64) string {
	return fmt.Sprintf("%.1f WPM", wpm)
}

// Format accuracy
func FormatAccuracy(accuracy float64) string {
	return fmt.Sprintf("%.1f%%", accuracy)
}
