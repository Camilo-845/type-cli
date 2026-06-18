package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content string

	switch m.screen {
	case screenMenu:
		content = m.viewMenu()
	case screenTyping:
		content = m.viewTyping()
	case screenResults:
		content = m.viewResults()
	case screenHistory:
		content = m.viewHistory()
	}

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
