package tui

import (
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

	wordPool := words.Generate(500, m.cfg.WordList, genCfg)

	if m.cfg.Mode == "words" {
		m.gm = game.NewWordGame(m.cfg.WordCount, wordPool)
	} else {
		m.gm = game.NewTimeGame(m.cfg.Duration, wordPool)
	}

	m.screen = screenTyping
	m.result = nil
}
