package tui

import (
	"strings"

	"github.com/Camilo-845/type-cli/internal/config"
	"github.com/Camilo-845/type-cli/internal/game"
	"github.com/Camilo-845/type-cli/internal/words"

	hg "github.com/Camilo-845/type-cli/internal/history"
)

type screen int

const (
	screenMenu screen = iota
	screenTyping
	screenResults
	screenHistory
)

type Model struct {
	screen        screen
	cfg           *config.Config
	gm            *game.Game
	result        *game.Result
	results       []game.Result
	cursor        int
	width         int
	height        int
	historyScroll int

	filtering    bool
	filterText   string
	filteredList []string
	filterCursor int
}

func NewModel() Model {
	cfg := config.Load()
	loaded, _ := hg.Load()

	return Model{
		screen:  screenMenu,
		cfg:     cfg,
		cursor:  0,
		width:   80,
		height:  24,
		results: loaded,
	}
}

func (m *Model) SetLanguage(lang string) {
	m.cfg.WordList = lang
	m.cfg.Validate()
}

func (m *Model) startGame() {
	_ = config.Save(m.cfg)

	genCfg := words.GeneratorConfig{
		Punctuation: m.cfg.Punctuation,
		LazyMode:    m.cfg.LazyMode,
		Numbers:     m.cfg.Numbers,
	}

	wordPool := words.Generate(game.WordPoolSize, m.cfg.WordList, genCfg)

	if m.cfg.Mode == "words" {
		m.gm = game.NewWordGame(m.cfg.WordCount, wordPool)
	} else {
		m.gm = game.NewTimeGame(m.cfg.Duration, wordPool)
	}

	m.screen = screenTyping
	m.result = nil
}

func (m *Model) enterFilter() {
	m.filtering = true
	m.filterText = ""
	m.filterCursor = 0
	m.updateFilteredList()
}

func (m *Model) exitFilter(apply bool) {
	if apply && m.filterText != "" && len(m.filteredList) > 0 && m.filterCursor < len(m.filteredList) {
		m.cfg.WordList = m.filteredList[m.filterCursor]
		m.cfg.Validate()
	}
	m.filtering = false
	m.filterText = ""
	m.filteredList = nil
	m.filterCursor = 0
}

func (m *Model) appendFilter(r rune) {
	m.filterText += string(r)
	m.filterCursor = 0
	m.updateFilteredList()
}

func (m *Model) backspaceFilter() {
	if len(m.filterText) > 0 {
		m.filterText = m.filterText[:len(m.filterText)-1]
		m.filterCursor = 0
		m.updateFilteredList()
	}
}

func (m *Model) updateFilteredList() {
	lower := strings.ToLower(m.filterText)
	m.filteredList = nil
	for _, name := range words.SortedNames() {
		if strings.Contains(strings.ToLower(name), lower) {
			m.filteredList = append(m.filteredList, name)
		}
	}
	if m.filterCursor >= len(m.filteredList) {
		m.filterCursor = 0
	}
}
