package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"type_game2/internal/game"
	hg "type_game2/internal/history"
)

func (m Model) handleTypingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.gm.Reset()
		m.gm = nil
		m.screen = screenMenu
		return m, nil

	default:
		m.gm.HandleKey(key)
	}

	if m.gm.State == game.Complete {
		return m.finishGame()
	}

	return m, nil
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if m.screen != screenTyping || m.gm == nil {
		return m, nil
	}

	m.gm.Tick()

	if m.gm.State == game.Complete {
		return m.finishGame()
	}

	return m, tickCmd()
}

func (m Model) finishGame() (tea.Model, tea.Cmd) {
	stats := m.gm.Stats()

	r := game.NewResult(
		stats,
		m.cfg.Mode,
		m.cfg.Duration,
		m.cfg.WordCount,
		m.cfg.WordList,
	)
	m.result = &r

	_ = hg.Append(r)

	m.screen = screenResults
	return m, nil
}
