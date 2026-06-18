package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q":
		return m, tea.Quit

	case "j", "down":
		m.cursor++
		if m.cursor > 2 {
			m.cursor = 0
		}

	case "k", "up":
		m.cursor--
		if m.cursor < 0 {
			m.cursor = 2
		}

	case "l", "right":
		m.cfg.Apply(m.cursor, true)

	case "h", "left":
		m.cfg.Apply(m.cursor, false)

	case " ", "tab":
		m.startGame()
		return m, tickCmd()
	}

	return m, nil
}
