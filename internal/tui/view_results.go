package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewResults() string {
	if m.result == nil {
		return ""
	}

	r := m.result

	title := titleStyle.Render("TEST COMPLETE")

	wpmLine := fmt.Sprintf("  %s %s",
		resultLabelStyle.Render("wpm"),
		resultWpmStyle.Render(fmt.Sprintf("%.0f", r.WPM)),
	)

	accLine := fmt.Sprintf("  %s %s",
		resultLabelStyle.Render("accuracy"),
		resultValueStyle.Render(fmt.Sprintf("%.1f%%", r.Accuracy)),
	)

	rawLine := fmt.Sprintf("  %s %s",
		resultLabelStyle.Render("raw"),
		resultValueStyle.Render(fmt.Sprintf("%.0f wpm", r.RawWPM)),
	)

	consLine := fmt.Sprintf("  %s %s",
		resultLabelStyle.Render("consistency"),
		resultValueStyle.Render(fmt.Sprintf("%.0f%%", r.Consistency)),
	)

	charsLine := fmt.Sprintf("  %s  %s%d  %s%d  %s%d",
		resultLabelStyle.Render("chars"),
		lipgloss.NewStyle().Foreground(green).Render("correct: "), r.CharsCorrect,
		lipgloss.NewStyle().Foreground(red).Render("incorrect: "), r.CharsIncorrect,
		lipgloss.NewStyle().Foreground(subtle).Render("extra: "), r.CharsExtra,
	)

	modeLine := fmt.Sprintf("  %s  %s  ·  %s",
		resultLabelStyle.Render("mode"),
		resultValueStyle.Render(r.ModeText()),
		resultValueStyle.Render(r.WordList),
	)

	help := helpStyle.Render("[enter] retry    [h] history    [esc] menu    [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			title,
			"",
			wpmLine,
			accLine,
			rawLine,
			consLine,
			"",
			charsLine,
			modeLine,
			"",
			help,
		),
	)
}
