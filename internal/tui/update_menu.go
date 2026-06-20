package tui

import (
	"unicode"

	tea "github.com/charmbracelet/bubbletea"

	hg "github.com/Camilo-845/type-cli/internal/history"
)

func (m Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.filtering {
		return m.handleFilterKey(msg)
	}

	maxCursor := m.cfg.CursorCount() - 1

	switch key {
	case "q":
		return m, tea.Quit

	case "j", "down":
		m.cursor++
		if m.cursor > maxCursor {
			m.cursor = 0
		}

	case "k", "up":
		m.cursor--
		if m.cursor < 0 {
			m.cursor = maxCursor
		}

	case "l", "right":
		if m.cursor == 2 {
			m.enterFilter()
		} else {
			m.cfg.Apply(m.cursor, true)
		}

	case "h", "left":
		m.cfg.Apply(m.cursor, false)

	case "H":
		m.previousScreen = m.screen
		m.screen = screenHistory
		m.historyScroll = 0
		var err error
		m.results, err = hg.Load()
		if err != nil {
			m.results = nil
		}
		return m, nil

	case " ", "enter":
		m.startGame()
		return m, tickCmd()

	default:
		if m.cursor == 2 && isPrintable(key) {
			m.enterFilter()
			m.appendFilter([]rune(key)[0])
		}
	}

	return m, nil
}

func (m Model) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.exitFilter(false)

	case "enter":
		m.exitFilter(true)

	case "backspace":
		m.backspaceFilter()

	case "down":
		if len(m.filteredList) > 0 {
			m.filterCursor++
			if m.filterCursor >= len(m.filteredList) {
				m.filterCursor = 0
			}
		}

	case "up":
		if len(m.filteredList) > 0 {
			m.filterCursor--
			if m.filterCursor < 0 {
				m.filterCursor = len(m.filteredList) - 1
			}
		}

	default:
		if isPrintable(key) {
			m.appendFilter([]rune(key)[0])
		}
	}

	return m, nil
}

func isPrintable(key string) bool {
	if len(key) != 1 {
		return false
	}
	r := rune(key[0])
	return unicode.IsPrint(r) && !unicode.IsSpace(r)
}
