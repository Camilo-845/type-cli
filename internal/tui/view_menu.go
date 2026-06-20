package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewMenu() string {
	logo := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		Render(strings.Join([]string{
			"╔╦╗╦ ╦╔═╗╔═╗   ╔═╗╦  ╦",
			" ║ ╚╦╝╠═╝║╣  ═ ║  ║  ║",
			" ╩  ╩ ╩  ╚═╝   ╚═╝╚═ ╩",
		}, "\n"))

	paramLabel := "duration"
	paramValue := fmt.Sprintf("%ds", m.cfg.Duration)
	if m.cfg.Mode == "words" {
		paramLabel = "word count"
		paramValue = fmt.Sprintf("%d", m.cfg.WordCount)
	}

	puncVal := "off"
	if m.cfg.Punctuation {
		puncVal = "on"
	}
	lazyVal := "off"
	if m.cfg.LazyMode {
		lazyVal = "on"
	}
	numsVal := "off"
	if m.cfg.Numbers {
		numsVal = "on"
	}

	items := []struct {
		label string
		value string
	}{
		{"mode", m.cfg.Mode},
		{paramLabel, paramValue},
		{"language", m.cfg.WordList},
		{"punctuation", puncVal},
		{"lazy mode", lazyVal},
		{"numbers", numsVal},
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

	helpStr := "[enter/space] start    [↑/↓] navigate    [←/→][h/l] change    [H] history    [q] quit"
	if m.width < 60 {
		helpStr = "[enter/space] start  [↑/↓] nav\n[h/l] change  [H] history  [q] quit"
	}
	helpLines := strings.Split(helpStr, "\n")
	padWidth := max(m.width, 30)
	for i, line := range helpLines {
		helpLines[i] = lipgloss.PlaceHorizontal(padWidth, lipgloss.Center, line)
	}
	help := helpStyle.Render(strings.Join(helpLines, "\n"))

	return renderContainer(m.width,
		lipgloss.JoinVertical(lipgloss.Center,
			logo,
			"",
			menu.String(),
			"",
			help,
		),
	)
}
