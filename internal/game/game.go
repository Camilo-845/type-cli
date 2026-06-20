package game

import (
	"time"
)

type State int

const (
	Idle State = iota
	Running
	Paused
	Complete
)

type WordState struct {
	Word    []rune
	Typed   []rune
	Correct []bool
}

type Game struct {
	Words   []WordState
	Current int
	State   State
	Elapsed time.Duration

	keystrokes   KeystrokeTracker
	bursts       BurstTracker
	firstKeyTime time.Time
	totalChars   int

	TimeMode  bool
	Duration  time.Duration
	WordCount int
}

type Stats struct {
	WPM             float64
	RawWPM          float64
	Accuracy        float64
	Consistency     float64
	CharsCorrect    int
	CharsIncorrect  int
	CharsExtra      int
	WordsTyped      int
	Duration        time.Duration
}

const WordPoolSize = 500

func NewTimeGame(durationSec int, words []string) *Game {
	batch := make([]WordState, WordPoolSize)
	var totalChars int
	for i, w := range words {
		if i >= WordPoolSize {
			break
		}
		batch[i] = WordState{Word: []rune(w), Typed: nil, Correct: nil}
		totalChars += len([]rune(w))
	}
	return &Game{
		Words:      batch,
		Current:    0,
		State:      Idle,
		TimeMode:   true,
		Duration:   time.Duration(durationSec) * time.Second,
		WordCount:  WordPoolSize,
		totalChars: totalChars,
	}
}

func NewWordGame(count int, words []string) *Game {
	batch := make([]WordState, count)
	var totalChars int
	for i, w := range words {
		if i >= count {
			break
		}
		batch[i] = WordState{Word: []rune(w), Typed: nil, Correct: nil}
		totalChars += len([]rune(w))
	}
	return &Game{
		Words:      batch,
		Current:    0,
		State:      Idle,
		TimeMode:   false,
		WordCount:  count,
		totalChars: totalChars,
	}
}

func (g *Game) CurrentWord() *WordState {
	if g.Current >= len(g.Words) {
		return nil
	}
	return &g.Words[g.Current]
}

func (g *Game) Remaining() time.Duration {
	if !g.TimeMode {
		return 0
	}
	remaining := g.Duration - g.Elapsed
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}
