package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	green     = lipgloss.Color("2")
	red       = lipgloss.Color("1")
	subtle    = lipgloss.Color("8")
	muted     = lipgloss.Color("7")
	highlight = lipgloss.Color("6")
	gold      = lipgloss.Color("3")
	white     = lipgloss.Color("15")
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
			MarginTop(1)
)

func renderContainer(width int, content string) string {
	pad := 2
	if width < 40 {
		pad = 1
	}
	if width < 30 {
		pad = 0
	}
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Padding(0, pad).
		Render(content)
}

var (
	correctStyle = lipgloss.NewStyle().
			Foreground(green)

	incorrectStyle = lipgloss.NewStyle().
			Foreground(red).
			Strikethrough(true)

	untypedStyle = lipgloss.NewStyle().
			Foreground(white)

	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Underline(true)
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
