package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"type_game2/internal/game"
	hg "type_game2/internal/history"
)

const tickInterval = time.Second

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		return m.handleTick()
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.screen {
	case screenMenu:
		return m.handleMenuKey(msg)
	case screenTyping:
		return m.handleTypingKey(msg)
	case screenResults:
		return m.handleResultKey(msg)
	case screenHistory:
		return m.handleHistoryKey(msg)
	}
	return m, nil
}

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

	case "enter", "right", "left":
		switch m.cursor {
		case 0:
			if m.cfg.Mode == "time" {
				m.cfg.Mode = "words"
			} else {
				m.cfg.Mode = "time"
			}
		case 1:
			if m.cfg.Mode == "time" {
				durations := []int{15, 30, 60, 120}
				for i, d := range durations {
					if d == m.cfg.Duration {
						m.cfg.Duration = durations[(i+1)%len(durations)]
						break
					}
				}
			} else {
				counts := []int{10, 25, 50, 100}
				for i, c := range counts {
					if c == m.cfg.WordCount {
						m.cfg.WordCount = counts[(i+1)%len(counts)]
						break
					}
				}
			}
		case 2:
			lists := []string{"english", "english_1k"}
			for i, l := range lists {
				if l == m.cfg.WordList {
					m.cfg.WordList = lists[(i+1)%len(lists)]
					break
				}
			}
		}

	case " ", "tab":
		m.startGame()
		return m, tickCmd()
	}

	return m, nil
}

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

func (m Model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q":
		return m, tea.Quit

	case "h":
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
	case "esc":
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
		maxScroll := len(m.results) - 15
		if maxScroll < 0 {
			maxScroll = 0
		}
		if m.historyScroll < maxScroll {
			m.historyScroll++
		}

	case "k", "up":
		if m.historyScroll > 0 {
			m.historyScroll--
		}

	case "esc":
		m.screen = screenResults
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
