package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"type_game2/internal/game"
)

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
