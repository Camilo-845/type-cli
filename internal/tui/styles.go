package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	green     = lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#4CAF50"}
	red       = lipgloss.AdaptiveColor{Light: "#C62828", Dark: "#F44336"}
	subtle    = lipgloss.AdaptiveColor{Light: "#9E9E9E", Dark: "#616161"}
	muted     = lipgloss.AdaptiveColor{Light: "#BDBDBD", Dark: "#424242"}
	highlight = lipgloss.AdaptiveColor{Light: "#1565C0", Dark: "#64B5F6"}
	gold      = lipgloss.AdaptiveColor{Light: "#F57F17", Dark: "#FFD54F"}
	white     = lipgloss.AdaptiveColor{Light: "#212121", Dark: "#E0E0E0"}
	bg        = lipgloss.AdaptiveColor{Light: "#FAFAFA", Dark: "#1E1E1E"}
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginBottom(2)

	menuOptionStyle = lipgloss.NewStyle().
			Foreground(white).
			Padding(0, 1)

	menuSelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Padding(0, 1)

	menuValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(gold)

	helpStyle = lipgloss.NewStyle().
			Foreground(muted).
			MarginTop(2)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(1, 2)
)

var (
	correctStyle = lipgloss.NewStyle().
			Foreground(green)

	incorrectStyle = lipgloss.NewStyle().
			Foreground(red).
			Strikethrough(true)

	untypedStyle = lipgloss.NewStyle().
			Foreground(subtle)

	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Underline(true)

	typedWordStyle = lipgloss.NewStyle().
			Foreground(muted)
)

var (
	resultWpmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(gold)

	resultLabelStyle = lipgloss.NewStyle().
			Foreground(subtle)

	resultValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)
)

var (
	historyHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Underline(true)

	historyItemStyle = lipgloss.NewStyle().
			Foreground(white)

	historySelectedStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)
)

var (
	topBarStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginBottom(1)

	timerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)
)
