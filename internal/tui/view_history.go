package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewHistory() string {
	title := titleStyle.Render("History")

	if len(m.results) == 0 {
		return renderContainer(m.width,
			lipgloss.JoinVertical(lipgloss.Center,
				title,
				"",
				subtitleStyle.Render("No tests completed yet"),
				"",
				helpStyle.Render("[tab/esc] back    [q] quit"),
			),
		)
	}

	narrow := m.width < 60

	var header string
	if narrow {
		header = historyHeaderStyle.Render(fmt.Sprintf("  %-8s  %-5s  %6s  %5s",
			"date", "mode", "wpm", "acc"))
	} else {
		header = historyHeaderStyle.Render(fmt.Sprintf("  %-12s  %-6s  %-12s  %8s  %6s",
			"date", "mode", "words", "wpm", "acc"))
	}

	total := len(m.results)
	newest := total - 1 - m.historyScroll
	oldest := max(newest-14, 0)

	var items []string
	items = append(items, header)
	items = append(items, "")

	for i := newest; i >= oldest; i-- {
		r := m.results[i]

		style := historyItemStyle
		if i == newest {
			style = historySelectedStyle
		}

		durationText := fmt.Sprintf("%ds", r.Duration)
		if r.Mode == "words" {
			durationText = fmt.Sprintf("%dw", r.WordCount)
		}

		wordList := r.WordList
		if narrow && len(wordList) > 8 {
			wordList = wordList[:8]
		}

		var line string
		if narrow {
			line = style.Render(fmt.Sprintf("  %-8s  %-5s  %6.0f  %5.1f%%",
				r.Date.Format("01-02"),
				durationText,
				r.WPM,
				r.Accuracy,
			))
		} else {
			line = style.Render(fmt.Sprintf("  %-12s  %-6s  %-12s  %6.0f  %5.1f%%",
				r.Date.Format("2006-01-02"),
				durationText,
				wordList,
				r.WPM,
				r.Accuracy,
			))
		}
		items = append(items, line)
	}

	help := helpStyle.Render("[↑/↓] scroll    [tab/esc] back    [q] quit")

	items = append(items, "")
	items = append(items, help)

	return renderContainer(m.width,
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}
