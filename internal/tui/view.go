package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"type_game2/internal/game"
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

func (m Model) viewMenu() string {
	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		Render(strings.Join([]string{
			"╔╦╗╦ ╦╔═╗╔═╗╔═╗╔╦╗╔═╗",
			" ║ ╚╦╝╠═╝║╣ ║ ╦║║║║╣ ",
			" ╩  ╩ ╩  ╚═╝╚═╝╩ ╩╚═╝",
		}, "\n"))

	paramLabel := "duration"
	paramValue := fmt.Sprintf("%ds", m.cfg.Duration)
	if m.cfg.Mode == "words" {
		paramLabel = "word count"
		paramValue = fmt.Sprintf("%d", m.cfg.WordCount)
	}

	items := []struct {
		label string
		value string
	}{
		{"mode", m.cfg.Mode},
		{paramLabel, paramValue},
		{"word list", m.cfg.WordList},
	}

	var menu strings.Builder
	for i, item := range items {
		labelStyle := menuOptionStyle
		prefix := "  "
		if i == m.cursor {
			labelStyle = menuSelectedStyle
			prefix = "▸ "
		}

		fmt.Fprintf(&menu, "%s%-14s%s\n",
			prefix,
			labelStyle.Render(item.label),
			menuValueStyle.Render("["+item.value+"]"))
	}

	help := helpStyle.Render("[space] start    [↑/↓] navigate    [←/→][h/l] change    [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			logo,
			"",
			menu.String(),
			"",
			help,
		),
	)
}

func (m Model) viewTyping() string {
	if m.gm == nil {
		return "loading..."
	}

	topBar := m.renderTopBar()

	wordsArea := m.renderWordsArea()

	progressArea := m.renderProgressArea()

	help := helpStyle.Render("[esc] back to menu")

	return lipgloss.JoinVertical(lipgloss.Center,
		topBar,
		"",
		wordsArea,
		"",
		progressArea,
		help,
	)
}

func (m Model) renderTopBar() string {
	modeText := fmt.Sprintf("%ds", m.cfg.Duration)
	if m.cfg.Mode == "words" {
		modeText = fmt.Sprintf("%d words", m.cfg.WordCount)
	}

	left := topBarStyle.Render(fmt.Sprintf("%s    %s", modeText, m.cfg.WordList))

	right := ""
	if m.gm.State == game.Running {
		right = topBarStyle.Render(fmt.Sprintf("%.0f wpm   %.0f%%",
			m.gm.LiveWPM(), m.gm.LiveAccuracy()))
	}

	barWidth := max(m.width-4, 40)

	leftPad := barWidth - lipgloss.Width(left) - lipgloss.Width(right)
	if leftPad < 0 {
		leftPad = 1
	}

	return left + strings.Repeat(" ", leftPad) + right
}

func (m Model) renderWordsArea() string {
	if m.gm == nil {
		return ""
	}

	maxWidth := max(m.width-8, 20)

	current := m.gm.Current
	start := max(current-3, 0)

	end := min(start+50, len(m.gm.Words))

	var lines []string
	var currentLineWords []string
	currentLineWidth := 0

	for i := start; i < end; i++ {
		ws := m.gm.Words[i]

		var word string
		if i < current {
			word = typedWordStyle.Render(ws.Word)
		} else if i == current {
			word = m.renderCurrentWord(ws)
		} else {
			word = untypedStyle.Render(ws.Word)
		}

		wordWidth := lipgloss.Width(word)

		if currentLineWidth+wordWidth+1 > maxWidth && len(currentLineWords) > 0 {
			lines = append(lines, strings.Join(currentLineWords, " "))
			currentLineWords = nil
			currentLineWidth = 0
		}

		if len(currentLineWords) > 0 {
			currentLineWidth++
		}
		currentLineWords = append(currentLineWords, word)
		currentLineWidth += wordWidth

		if len(lines) >= 3 {
			break
		}
	}

	if len(currentLineWords) > 0 {
		lines = append(lines, strings.Join(currentLineWords, " "))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderCurrentWord(ws game.WordState) string {
	var chars []string

	for i, ch := range ws.Word {
		switch {
		case i < len(ws.Correct) && ws.Correct[i]:
			chars = append(chars, correctStyle.Render(string(ch)))
		case i < len(ws.Correct) && !ws.Correct[i]:
			chars = append(chars, incorrectStyle.Render(string(ch)))
		case i == len(ws.Typed) && (m.gm.State == game.Running || m.gm.State == game.Idle):
			chars = append(chars, cursorStyle.Render(string(ch)))
		default:
			chars = append(chars, lipgloss.NewStyle().Foreground(white).Render(string(ch)))
		}
	}

	for i := len(ws.Word); i < len(ws.Typed); i++ {
		chars = append(chars, incorrectStyle.Render(string(ws.Typed[i])))
	}

	return strings.Join(chars, "")
}

func (m Model) renderProgressArea() string {
	if m.gm == nil {
		return ""
	}

	barWidth := 30
	if m.width > 0 {
		barWidth = min(max(m.width-20, 10), 60)
	}

	var filled int
	var label string

	totalChars := m.gm.TotalChars()
	if totalChars > 0 {
		filled = int(float64(m.gm.CorrectChars()) / float64(totalChars) * float64(barWidth))
	}

	if m.gm.TimeMode {
		remaining := m.gm.Remaining()
		label = fmt.Sprintf("%.0fs", remaining.Seconds())
	} else {
		label = fmt.Sprintf("%d/%d", m.gm.Current, m.gm.WordCount)
	}

	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	bar := ""
	if m.gm.State == game.Running || m.gm.State == game.Idle {
		bar = lipgloss.NewStyle().Foreground(gold).Render(strings.Repeat("▁", filled))
		bar += lipgloss.NewStyle().Foreground(muted).Render(strings.Repeat("▁", barWidth-filled))
	} else {
		bar = lipgloss.NewStyle().Foreground(muted).Render(strings.Repeat("▁", barWidth))
	}

	labelStyled := lipgloss.NewStyle().Foreground(subtle).Render(label)

	return fmt.Sprintf("%s  %s", bar, labelStyled)
}

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

func (m Model) viewHistory() string {
	title := titleStyle.Render("History")

	if len(m.results) == 0 {
		return containerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Center,
				title,
				"",
				subtitleStyle.Render("No tests completed yet"),
				"",
				helpStyle.Render("[esc] back    [q] quit"),
			),
		)
	}

	header := historyHeaderStyle.Render(fmt.Sprintf("  %-12s  %-6s  %-12s  %8s  %6s",
		"date", "mode", "words", "wpm", "acc"))

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

		line := style.Render(fmt.Sprintf("  %-12s  %-6s  %-12s  %6.0f  %5.1f%%",
			r.Date.Format("2006-01-02"),
			durationText,
			r.WordList,
			r.WPM,
			r.Accuracy,
		))
		items = append(items, line)
	}

	help := helpStyle.Render("[↑/↓] scroll    [esc] back    [q] quit")

	items = append(items, "")
	items = append(items, help)

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}
