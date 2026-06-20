package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	hg 	"github.com/Camilo-845/type-cli/internal/history"
)

const historyVisibleEntries = 15

func (m Model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q":
		return m, tea.Quit

	case "h":
		m.previousScreen = m.screen
		m.screen = screenHistory
		m.historyScroll = 0
		var err error
		m.results, err = hg.Load()
		if err != nil {
			m.results = nil
		}
		return m, nil

	case "r", "enter", " ":
		fallthrough
	case "esc", "tab":
		m.screen = screenMenu
		m.result = nil
		m.gm = nil
		return m, nil
	}

	return m, nil
}

func (m Model) handleHistoryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q":
		return m, tea.Quit

	case "j", "down":
		maxScroll := max(len(m.results)-historyVisibleEntries, 0)
		if m.historyScroll < maxScroll {
			m.historyScroll++
		}

	case "k", "up":
		if m.historyScroll > 0 {
			m.historyScroll--
		}

	case "esc":
		m.screen = m.previousScreen
		m.historyScroll = 0
		return m, nil

	case "r", "enter", " ":
		m.screen = screenMenu
		m.result = nil
		m.gm = nil
		return m, nil
	}

	return m, nil
}
