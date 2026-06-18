package tui

import (
	"github.com/Camilo-845/typingame/internal/config"
	"github.com/Camilo-845/typingame/internal/game"
	"github.com/Camilo-845/typingame/internal/words"

	hg "github.com/Camilo-845/typingame/internal/history"
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

func (m *Model) startGame() {
	_ = config.Save(m.cfg)

	wordPool := words.Generate(500, m.cfg.WordList)

	if m.cfg.Mode == "words" {
		m.gm = game.NewWordGame(m.cfg.WordCount, wordPool)
	} else {
		m.gm = game.NewTimeGame(m.cfg.Duration, wordPool)
	}

	m.screen = screenTyping
	m.result = nil
}
